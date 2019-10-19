var dateFormat = require('dateformat');
const {grpc} = require("@improbable-eng/grpc-web");
const {Empty} = require("google-protobuf/google/protobuf/empty_pb")

const {TemperatureServiceClient} = require('../proto/temperature_pb_service.js')

const client = new TemperatureServiceClient(window.env.API);
const empty = new Empty()

response = client.getTemperature(empty)
response.on("data", function(message) {
  temperature = message.toObject().temperature
  hue = 180 - (180 * ((temperature + 20) / 70))
  document.getElementById("temperature").innerHTML = temperature + " °C";
  document.getElementById("box").style.background = "hsl(" + hue + ", 100%, 50%)";

  updated = dateFormat("mmm dS, yyyy @ HH:MM:ss")
  document.getElementById("last-updated").innerHTML = updated;
}).on("status", function(data) {
  console.log(data)
}).on("end", function(data) {
  console.log(data)
})

