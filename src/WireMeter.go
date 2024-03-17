package main

import (
	//_ "net/http/pprof"

	"archive/zip"
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"log"
	"net/http"
)

//go:embed views/*
var viewsFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {

	/*go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()*/

	var receivingInterface string
	var sendingInterface string
	var waitTime int
	flag.StringVar(&receivingInterface, "r", "enp1", "receiving interface")
	flag.StringVar(&sendingInterface, "s", "enp2", "sending interface")
	flag.IntVar(&waitTime, "sleep", 5, "time to wait in between send requests (good values may depend on your network interface)")
	flag.Parse()
	/* TODO: add check to ensure selected interfaces exist */

	sensor := NewWireSensor(receivingInterface, sendingInterface, 1024, waitTime)
	sensor.run()

	engine := html.NewFileSystem(http.FS(viewsFS), ".html")
	app := fiber.New(fiber.Config{
		Views:        engine,
		ServerHeader: "WireMeter",
		AppName:      "WireMeter",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// TODO: add logging here
			fmt.Println(err)
			return c.Render("views/error", fiber.Map{}, "views/layouts/main")

		},
	})
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(staticFS),
		PathPrefix: "static",
		Browse:     true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		env := fiber.Map{}
		return c.Render("views/index", env, "views/layouts/main")
	})

	app.Get("api/measurements", func(c *fiber.Ctx) error {

		type measurement struct {
			Timestamps      []string `json:"timestamps"`
			AverageDuration []int64  `json:"averagedur"`
			MinDuration     []int64  `json:"mindur"`
			MaxDuration     []int64  `json:"maxdur"`
			ResolvedPackets []int64  `json:"resolvedpkts"`
			UnknownPackets  []int64  `json:"unknownpkts"`
			AgedPackets     []int64  `json:"agedpkts"`
		}

		timestamps, averageDurations, minDurations, maxDurations, resolvedPackets, unknownPackets, agedPackets := sensor.measurements.analyzeMeasurements()

		m := measurement{
			Timestamps:      timestamps,
			AverageDuration: averageDurations,
			MinDuration:     minDurations,
			MaxDuration:     maxDurations,
			ResolvedPackets: resolvedPackets,
			UnknownPackets:  unknownPackets,
			AgedPackets:     agedPackets,
		}

		foo_marshalled, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return c.SendString(string(foo_marshalled))
	})

	app.Get("api/measurements/snapshot", func(c *fiber.Ctx) error {
		c.Set("Content-Disposition", "attachment; filename=snapshot.svg")
		c.Set("Content-Type", "image/svg+xml")
		runtime, loss := sensor.measurements.exportSVG()
		fileContents := map[string]string{
			"runtime.svg": runtime,
			"loss.svg":    loss,
			"raw.csv":     sensor.measurements.exportCSV(),
		}
		zipBuffer, err := createZipFile(fileContents)
		if err != nil {
			log.Fatal(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		c.Set(fiber.HeaderContentType, "application/zip")
		c.Set(fiber.HeaderContentDisposition, "attachment; filename=snapshot.zip")

		// Return the ZIP file as a response
		return c.Send(zipBuffer.Bytes())
	})

	log.Fatal(app.Listen(":3000"))
}

func createZipFile(fileContents map[string]string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for fileName, content := range fileContents {
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return nil, err
		}

		_, err = fileWriter.Write([]byte(content))
		if err != nil {
			return nil, err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
