// Set new default font family and font color to mimic Bootstrap's default styling
Chart.defaults.global.defaultFontFamily = 'Nunito', '-apple-system,system-ui,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,sans-serif';
Chart.defaults.global.defaultFontColor = '#858796';

function number_format(number, decimals, dec_point, thousands_sep) {
    // *     example: number_format(1234.56, 2, ',', ' ');
    // *     return: '1 234,56'
    number = (number + '').replace(',', '').replace(' ', '');
    var n = !isFinite(+number) ? 0 : +number,
        prec = !isFinite(+decimals) ? 0 : Math.abs(decimals),
        sep = (typeof thousands_sep === 'undefined') ? ',' : thousands_sep,
        dec = (typeof dec_point === 'undefined') ? '.' : dec_point,
        s = '',
        toFixedFix = function(n, prec) {
            var k = Math.pow(10, prec);
            return '' + Math.round(n * k) / k;
        };
    // Fix for IE parseFloat(0.55).toFixed(0) = 0;
    s = (prec ? toFixedFix(n, prec) : '' + Math.round(n)).split('.');
    if (s[0].length > 3) {
        s[0] = s[0].replace(/\B(?=(?:\d{3})+(?!\d))/g, sep);
    }
    if ((s[1] || '').length < prec) {
        s[1] = s[1] || '';
        s[1] += new Array(prec - s[1].length + 1).join('0');
    }
    return s.join(dec);
}

// Area Chart Example
var runtimeChart = document.getElementById("RuntimeChart");
var runtimeChartSwitch = document.getElementById("RuntimeChartSwitch");
var typeChart = document.getElementById("TypeChart");
var typeChartSwitch = document.getElementById("TypeChartSwitch");
var color_blue = "rgba(78, 115, 223, 1)"
var color_red = "rgba(231, 74, 59, 1)"
var color_green = "rgba(28, 200, 138, 1)"

var myLineChart = new Chart(runtimeChart, {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: "Average",
            lineTension: 0.3,
            backgroundColor: "rgba(78, 115, 223, 0.05)",
            borderColor: color_blue,
            pointRadius: 3,
            pointBackgroundColor: color_blue,
            pointBorderColor: color_blue,
            pointHoverRadius: 3,
            pointHoverBackgroundColor: color_blue,
            pointHoverBackgroundColor: color_blue,
            pointHoverBorderColor: color_blue,
            pointHitRadius: 10,
            pointBorderWidth: 2,
            data: [],
        },
            {
                label: "Min",
                lineTension: 0.3,
                backgroundColor: "rgba(28, 200, 138, 0.05)",
                borderColor: color_green,
                pointRadius: 3,
                pointBackgroundColor: color_green,
                pointBorderColor: color_green,
                pointHoverRadius: 3,
                pointHoverBackgroundColor: color_green,
                pointHoverBorderColor: color_green,
                pointHitRadius: 10,
                pointBorderWidth: 2,
                data: [],
            },
            {
                label: "Max",
                lineTension: 0.3,
                backgroundColor: "rgba(231, 74, 59, 0.05)",
                borderColor: color_red,
                pointRadius: 3,
                pointBackgroundColor: color_red,
                pointBorderColor: color_red,
                pointHoverRadius: 3,
                pointHoverBackgroundColor: color_red,
                pointHoverBorderColor: color_red,
                pointHitRadius: 10,
                pointBorderWidth: 2,
                data: [],
            }],
    },
    options: {
        maintainAspectRatio: false,
        layout: {
            padding: {
                left: 10,
                right: 25,
                top: 25,
                bottom: 0
            }
        },
        scales: {
            xAxes: [{
                time: {
                    unit: 'date'
                },
                gridLines: {
                    display: false,
                    drawBorder: false
                },
                ticks: {
                    maxTicksLimit: 7
                }
            }],
            yAxes: [{
                ticks: {
                    maxTicksLimit: 5,
                    padding: 10,
                    // Include a dollar sign in the ticks
                    callback: function(value, index, values) {
                        return number_format(value)+' ms';
                    }
                },
                gridLines: {
                    color: "rgb(234, 236, 244)",
                    zeroLineColor: "rgb(234, 236, 244)",
                    drawBorder: false,
                    borderDash: [2],
                    zeroLineBorderDash: [2]
                }
            }],
        },
        legend: {
            display: true
        },
        tooltips: {
            backgroundColor: "rgb(255,255,255)",
            bodyFontColor: "#858796",
            titleMarginBottom: 10,
            titleFontColor: '#6e707e',
            titleFontSize: 14,
            borderColor: '#dddfeb',
            borderWidth: 1,
            xPadding: 15,
            yPadding: 15,
            displayColors: false,
            intersect: false,
            mode: 'index',
            caretPadding: 10,
            callbacks: {
                label: function(tooltipItem, chart) {
                    var datasetLabel = chart.datasets[tooltipItem.datasetIndex].label || '';
                    return datasetLabel + ': ' + number_format(tooltipItem.yLabel) + ' ms';
                }
            }
        }
    }
});
var myTypeChart = new Chart(typeChart, {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: "Resolved",
            lineTension: 0.3,
            backgroundColor: "rgba(78, 115, 223, 0.05)",
            borderColor: color_blue,
            pointRadius: 3,
            pointBackgroundColor: color_blue,
            pointBorderColor: color_blue,
            pointHoverRadius: 3,
            pointHoverBackgroundColor: color_blue,
            pointHoverBorderColor: color_blue,
            pointHitRadius: 10,
            pointBorderWidth: 2,
            data: [],
        },
            {
                label: "Unknown",
                lineTension: 0.3,
                backgroundColor: "rgba(28, 200, 138, 0.05)",
                borderColor: color_green,
                pointRadius: 3,
                pointBackgroundColor: color_green,
                pointBorderColor: color_green,
                pointHoverRadius: 3,
                pointHoverBackgroundColor: color_green,
                pointHoverBorderColor: color_green,
                pointHitRadius: 10,
                pointBorderWidth: 2,
                data: [],
            },
            {
                label: "Aged",
                lineTension: 0.3,
                backgroundColor: "rgba(231, 74, 59, 0.05)",
                borderColor: color_red,
                pointRadius: 3,
                pointBackgroundColor: color_red,
                pointBorderColor: color_red,
                pointHoverRadius: 3,
                pointHoverBackgroundColor: color_red,
                pointHoverBorderColor: color_red,
                pointHitRadius: 10,
                pointBorderWidth: 2,
                data: [],
            }],
    },
    options: {
        maintainAspectRatio: false,
        layout: {
            padding: {
                left: 10,
                right: 25,
                top: 25,
                bottom: 0
            }
        },
        scales: {
            xAxes: [{
                time: {
                    unit: 'date'
                },
                gridLines: {
                    display: false,
                    drawBorder: false
                },
                ticks: {
                    maxTicksLimit: 7
                }
            }],
            yAxes: [{
                ticks: {
                    maxTicksLimit: 5,
                    padding: 10,
                    // Include a dollar sign in the ticks
                    callback: function(value, index, values) {
                        return number_format(value)+' packets';
                    }
                },
                gridLines: {
                    color: "rgb(234, 236, 244)",
                    zeroLineColor: "rgb(234, 236, 244)",
                    drawBorder: false,
                    borderDash: [2],
                    zeroLineBorderDash: [2]
                }
            }],
        },
        legend: {
            display: true
        },
        tooltips: {
            backgroundColor: "rgb(255,255,255)",
            bodyFontColor: "#858796",
            titleMarginBottom: 10,
            titleFontColor: '#6e707e',
            titleFontSize: 14,
            borderColor: '#dddfeb',
            borderWidth: 1,
            xPadding: 15,
            yPadding: 15,
            displayColors: false,
            intersect: false,
            mode: 'index',
            caretPadding: 10,
            callbacks: {
                label: function(tooltipItem, chart) {
                    var datasetLabel = chart.datasets[tooltipItem.datasetIndex].label || '';
                    return datasetLabel + ': ' + number_format(tooltipItem.yLabel) + ' packets';
                }
            }
        }
    }
});

var t=setInterval(refreshChart,1000);
async function refreshChart() {
    var dataset = {};
    console.log(runtimeChartSwitch)
    if ((!runtimeChartSwitch.checked)||(!typeChartSwitch.checked)){
        console.log("request data")
        const response = await fetch("/api/measurements");
        dataset = await response.json();
    }
    if (!runtimeChartSwitch.checked) {
        console.log("reprint runtime chart")
        myLineChart.data.datasets[0].data = dataset.averagedur
        myLineChart.data.datasets[1].data = dataset.mindur
        myLineChart.data.datasets[2].data = dataset.maxdur
        myLineChart.data.labels = dataset.timestamps
        myLineChart.update()
    }
    if (!typeChartSwitch.checked) {
        console.log("reprint type chart")
        myTypeChart.data.datasets[0].data = dataset.resolvedpkts
        myTypeChart.data.datasets[1].data = dataset.unknownpkts
        myTypeChart.data.datasets[2].data = dataset.agedpkts
        myTypeChart.data.labels = dataset.timestamps
        myTypeChart.update()
    }
}