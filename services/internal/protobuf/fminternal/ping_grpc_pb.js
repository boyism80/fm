// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var fminternal_ping_pb = require('../fminternal/ping_pb.js');

function serialize_fm_internal_CheckCharacterNameReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.CheckCharacterNameReply)) {
    throw new Error('Expected argument of type fm.internal.CheckCharacterNameReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CheckCharacterNameReply(buffer_arg) {
  return fminternal_ping_pb.CheckCharacterNameReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CheckCharacterNameRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.CheckCharacterNameRequest)) {
    throw new Error('Expected argument of type fm.internal.CheckCharacterNameRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CheckCharacterNameRequest(buffer_arg) {
  return fminternal_ping_pb.CheckCharacterNameRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CreateCharacterReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.CreateCharacterReply)) {
    throw new Error('Expected argument of type fm.internal.CreateCharacterReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CreateCharacterReply(buffer_arg) {
  return fminternal_ping_pb.CreateCharacterReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CreateCharacterRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.CreateCharacterRequest)) {
    throw new Error('Expected argument of type fm.internal.CreateCharacterRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CreateCharacterRequest(buffer_arg) {
  return fminternal_ping_pb.CreateCharacterRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DeleteCharacterReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.DeleteCharacterReply)) {
    throw new Error('Expected argument of type fm.internal.DeleteCharacterReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DeleteCharacterReply(buffer_arg) {
  return fminternal_ping_pb.DeleteCharacterReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DeleteCharacterRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.DeleteCharacterRequest)) {
    throw new Error('Expected argument of type fm.internal.DeleteCharacterRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DeleteCharacterRequest(buffer_arg) {
  return fminternal_ping_pb.DeleteCharacterRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DeleteInventoryReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.DeleteInventoryReply)) {
    throw new Error('Expected argument of type fm.internal.DeleteInventoryReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DeleteInventoryReply(buffer_arg) {
  return fminternal_ping_pb.DeleteInventoryReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DeleteInventoryRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.DeleteInventoryRequest)) {
    throw new Error('Expected argument of type fm.internal.DeleteInventoryRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DeleteInventoryRequest(buffer_arg) {
  return fminternal_ping_pb.DeleteInventoryRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetCharacterListReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.GetCharacterListReply)) {
    throw new Error('Expected argument of type fm.internal.GetCharacterListReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetCharacterListReply(buffer_arg) {
  return fminternal_ping_pb.GetCharacterListReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetCharacterListRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.GetCharacterListRequest)) {
    throw new Error('Expected argument of type fm.internal.GetCharacterListRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetCharacterListRequest(buffer_arg) {
  return fminternal_ping_pb.GetCharacterListRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetCharacterReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.GetCharacterReply)) {
    throw new Error('Expected argument of type fm.internal.GetCharacterReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetCharacterReply(buffer_arg) {
  return fminternal_ping_pb.GetCharacterReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetCharacterRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.GetCharacterRequest)) {
    throw new Error('Expected argument of type fm.internal.GetCharacterRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetCharacterRequest(buffer_arg) {
  return fminternal_ping_pb.GetCharacterRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetInventoryReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.GetInventoryReply)) {
    throw new Error('Expected argument of type fm.internal.GetInventoryReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetInventoryReply(buffer_arg) {
  return fminternal_ping_pb.GetInventoryReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetInventoryRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.GetInventoryRequest)) {
    throw new Error('Expected argument of type fm.internal.GetInventoryRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetInventoryRequest(buffer_arg) {
  return fminternal_ping_pb.GetInventoryRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LoginAccountReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.LoginAccountReply)) {
    throw new Error('Expected argument of type fm.internal.LoginAccountReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LoginAccountReply(buffer_arg) {
  return fminternal_ping_pb.LoginAccountReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LoginAccountRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.LoginAccountRequest)) {
    throw new Error('Expected argument of type fm.internal.LoginAccountRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LoginAccountRequest(buffer_arg) {
  return fminternal_ping_pb.LoginAccountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

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

function serialize_fm_internal_SaveCharacterReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.SaveCharacterReply)) {
    throw new Error('Expected argument of type fm.internal.SaveCharacterReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharacterReply(buffer_arg) {
  return fminternal_ping_pb.SaveCharacterReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveCharacterRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.SaveCharacterRequest)) {
    throw new Error('Expected argument of type fm.internal.SaveCharacterRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharacterRequest(buffer_arg) {
  return fminternal_ping_pb.SaveCharacterRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveCharactersReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.SaveCharactersReply)) {
    throw new Error('Expected argument of type fm.internal.SaveCharactersReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharactersReply(buffer_arg) {
  return fminternal_ping_pb.SaveCharactersReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveCharactersRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.SaveCharactersRequest)) {
    throw new Error('Expected argument of type fm.internal.SaveCharactersRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharactersRequest(buffer_arg) {
  return fminternal_ping_pb.SaveCharactersRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveInventoryReply(arg) {
  if (!(arg instanceof fminternal_ping_pb.SaveInventoryReply)) {
    throw new Error('Expected argument of type fm.internal.SaveInventoryReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveInventoryReply(buffer_arg) {
  return fminternal_ping_pb.SaveInventoryReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveInventoryRequest(arg) {
  if (!(arg instanceof fminternal_ping_pb.SaveInventoryRequest)) {
    throw new Error('Expected argument of type fm.internal.SaveInventoryRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveInventoryRequest(buffer_arg) {
  return fminternal_ping_pb.SaveInventoryRequest.deserializeBinary(new Uint8Array(buffer_arg));
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
  getCharacter: {
    path: '/fm.internal.Internal/GetCharacter',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.GetCharacterRequest,
    responseType: fminternal_ping_pb.GetCharacterReply,
    requestSerialize: serialize_fm_internal_GetCharacterRequest,
    requestDeserialize: deserialize_fm_internal_GetCharacterRequest,
    responseSerialize: serialize_fm_internal_GetCharacterReply,
    responseDeserialize: deserialize_fm_internal_GetCharacterReply,
  },
  saveCharacter: {
    path: '/fm.internal.Internal/SaveCharacter',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.SaveCharacterRequest,
    responseType: fminternal_ping_pb.SaveCharacterReply,
    requestSerialize: serialize_fm_internal_SaveCharacterRequest,
    requestDeserialize: deserialize_fm_internal_SaveCharacterRequest,
    responseSerialize: serialize_fm_internal_SaveCharacterReply,
    responseDeserialize: deserialize_fm_internal_SaveCharacterReply,
  },
  saveCharacters: {
    path: '/fm.internal.Internal/SaveCharacters',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.SaveCharactersRequest,
    responseType: fminternal_ping_pb.SaveCharactersReply,
    requestSerialize: serialize_fm_internal_SaveCharactersRequest,
    requestDeserialize: deserialize_fm_internal_SaveCharactersRequest,
    responseSerialize: serialize_fm_internal_SaveCharactersReply,
    responseDeserialize: deserialize_fm_internal_SaveCharactersReply,
  },
  getInventory: {
    path: '/fm.internal.Internal/GetInventory',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.GetInventoryRequest,
    responseType: fminternal_ping_pb.GetInventoryReply,
    requestSerialize: serialize_fm_internal_GetInventoryRequest,
    requestDeserialize: deserialize_fm_internal_GetInventoryRequest,
    responseSerialize: serialize_fm_internal_GetInventoryReply,
    responseDeserialize: deserialize_fm_internal_GetInventoryReply,
  },
  saveInventory: {
    path: '/fm.internal.Internal/SaveInventory',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.SaveInventoryRequest,
    responseType: fminternal_ping_pb.SaveInventoryReply,
    requestSerialize: serialize_fm_internal_SaveInventoryRequest,
    requestDeserialize: deserialize_fm_internal_SaveInventoryRequest,
    responseSerialize: serialize_fm_internal_SaveInventoryReply,
    responseDeserialize: deserialize_fm_internal_SaveInventoryReply,
  },
  deleteInventory: {
    path: '/fm.internal.Internal/DeleteInventory',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.DeleteInventoryRequest,
    responseType: fminternal_ping_pb.DeleteInventoryReply,
    requestSerialize: serialize_fm_internal_DeleteInventoryRequest,
    requestDeserialize: deserialize_fm_internal_DeleteInventoryRequest,
    responseSerialize: serialize_fm_internal_DeleteInventoryReply,
    responseDeserialize: deserialize_fm_internal_DeleteInventoryReply,
  },
  loginAccount: {
    path: '/fm.internal.Internal/LoginAccount',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.LoginAccountRequest,
    responseType: fminternal_ping_pb.LoginAccountReply,
    requestSerialize: serialize_fm_internal_LoginAccountRequest,
    requestDeserialize: deserialize_fm_internal_LoginAccountRequest,
    responseSerialize: serialize_fm_internal_LoginAccountReply,
    responseDeserialize: deserialize_fm_internal_LoginAccountReply,
  },
  getCharacterList: {
    path: '/fm.internal.Internal/GetCharacterList',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.GetCharacterListRequest,
    responseType: fminternal_ping_pb.GetCharacterListReply,
    requestSerialize: serialize_fm_internal_GetCharacterListRequest,
    requestDeserialize: deserialize_fm_internal_GetCharacterListRequest,
    responseSerialize: serialize_fm_internal_GetCharacterListReply,
    responseDeserialize: deserialize_fm_internal_GetCharacterListReply,
  },
  checkCharacterName: {
    path: '/fm.internal.Internal/CheckCharacterName',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.CheckCharacterNameRequest,
    responseType: fminternal_ping_pb.CheckCharacterNameReply,
    requestSerialize: serialize_fm_internal_CheckCharacterNameRequest,
    requestDeserialize: deserialize_fm_internal_CheckCharacterNameRequest,
    responseSerialize: serialize_fm_internal_CheckCharacterNameReply,
    responseDeserialize: deserialize_fm_internal_CheckCharacterNameReply,
  },
  createCharacter: {
    path: '/fm.internal.Internal/CreateCharacter',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.CreateCharacterRequest,
    responseType: fminternal_ping_pb.CreateCharacterReply,
    requestSerialize: serialize_fm_internal_CreateCharacterRequest,
    requestDeserialize: deserialize_fm_internal_CreateCharacterRequest,
    responseSerialize: serialize_fm_internal_CreateCharacterReply,
    responseDeserialize: deserialize_fm_internal_CreateCharacterReply,
  },
  deleteCharacter: {
    path: '/fm.internal.Internal/DeleteCharacter',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_ping_pb.DeleteCharacterRequest,
    responseType: fminternal_ping_pb.DeleteCharacterReply,
    requestSerialize: serialize_fm_internal_DeleteCharacterRequest,
    requestDeserialize: deserialize_fm_internal_DeleteCharacterRequest,
    responseSerialize: serialize_fm_internal_DeleteCharacterReply,
    responseDeserialize: deserialize_fm_internal_DeleteCharacterReply,
  },
};

exports.InternalClient = grpc.makeGenericClientConstructor(InternalService, 'Internal');
