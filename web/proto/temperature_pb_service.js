// package: api
// file: temperature.proto

var temperature_pb = require("./temperature_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var TemperatureService = (function () {
  function TemperatureService() {}
  TemperatureService.serviceName = "api.TemperatureService";
  return TemperatureService;
}());

TemperatureService.GetTemperature = {
  methodName: "GetTemperature",
  service: TemperatureService,
  requestStream: false,
  responseStream: true,
  requestType: google_protobuf_empty_pb.Empty,
  responseType: temperature_pb.TemperatureResponse
};

exports.TemperatureService = TemperatureService;

function TemperatureServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

TemperatureServiceClient.prototype.getTemperature = function getTemperature(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(TemperatureService.GetTemperature, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners.end.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

exports.TemperatureServiceClient = TemperatureServiceClient;

