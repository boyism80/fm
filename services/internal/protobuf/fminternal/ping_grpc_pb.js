// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var fminternal_ping_pb = require('../fminternal/ping_pb.js');

function serialize_fm_internal_PingReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.PingReply)) {
    throw new Error('Expected argument of type fm.internal.PingReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_PingReply(buffer_arg) {
  return fminternal_ping_pb.PingReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_PingRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.PingRequest)) {
    throw new Error('Expected argument of type fm.internal.PingRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_PingRequest(buffer_arg) {
  return fminternal_ping_pb.PingRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var InternalService = exports.InternalService = {
  ping: {
    path: '/fm.internal.Internal/Ping',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.PingRequest,
    responseType: fminternal_ping_pb.PingReply,
    requestSerialize: serialize_fm_internal_PingRequest,
    requestDeserialize: deserialize_fm_internal_PingRequest,
    responseSerialize: serialize_fm_internal_PingReply,
    responseDeserialize: deserialize_fm_internal_PingReply,
  },
};

exports.InternalClient = grpc.makeGenericClientConstructor(InternalService, 'Internal');
