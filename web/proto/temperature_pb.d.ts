// package: api
// file: temperature.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

export class TemperatureResponse extends jspb.Message {
  getTemperature(): number;
  setTemperature(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemperatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TemperatureResponse): TemperatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TemperatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemperatureResponse;
  static deserializeBinaryFromReader(message: TemperatureResponse, reader: jspb.BinaryReader): TemperatureResponse;
}

export namespace TemperatureResponse {
  export type AsObject = {
    temperature: number,
  }
}

