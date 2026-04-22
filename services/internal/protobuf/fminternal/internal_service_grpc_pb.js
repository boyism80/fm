// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var fminternal_internal_service_pb = require('../fminternal/internal_service_pb.js');

function serialize_fm_internal_BeginGameTransitionReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.BeginGameTransitionReply)) {
    throw new Error('Expected argument of type fm.internal.BeginGameTransitionReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_BeginGameTransitionReply(buffer_arg) {
  return fminternal_internal_service_pb.BeginGameTransitionReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_BeginGameTransitionRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.BeginGameTransitionRequest)) {
    throw new Error('Expected argument of type fm.internal.BeginGameTransitionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_BeginGameTransitionRequest(buffer_arg) {
  return fminternal_internal_service_pb.BeginGameTransitionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_BroadcastPartyChatReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.BroadcastPartyChatReply)) {
    throw new Error('Expected argument of type fm.internal.BroadcastPartyChatReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_BroadcastPartyChatReply(buffer_arg) {
  return fminternal_internal_service_pb.BroadcastPartyChatReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_BroadcastPartyChatRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.BroadcastPartyChatRequest)) {
    throw new Error('Expected argument of type fm.internal.BroadcastPartyChatRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_BroadcastPartyChatRequest(buffer_arg) {
  return fminternal_internal_service_pb.BroadcastPartyChatRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_ChangePartyLeaderReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.ChangePartyLeaderReply)) {
    throw new Error('Expected argument of type fm.internal.ChangePartyLeaderReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_ChangePartyLeaderReply(buffer_arg) {
  return fminternal_internal_service_pb.ChangePartyLeaderReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_ChangePartyLeaderRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.ChangePartyLeaderRequest)) {
    throw new Error('Expected argument of type fm.internal.ChangePartyLeaderRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_ChangePartyLeaderRequest(buffer_arg) {
  return fminternal_internal_service_pb.ChangePartyLeaderRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CheckCharacterNameReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.CheckCharacterNameReply)) {
    throw new Error('Expected argument of type fm.internal.CheckCharacterNameReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CheckCharacterNameReply(buffer_arg) {
  return fminternal_internal_service_pb.CheckCharacterNameReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CheckCharacterNameRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.CheckCharacterNameRequest)) {
    throw new Error('Expected argument of type fm.internal.CheckCharacterNameRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CheckCharacterNameRequest(buffer_arg) {
  return fminternal_internal_service_pb.CheckCharacterNameRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CreateCharacterReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.CreateCharacterReply)) {
    throw new Error('Expected argument of type fm.internal.CreateCharacterReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CreateCharacterReply(buffer_arg) {
  return fminternal_internal_service_pb.CreateCharacterReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CreateCharacterRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.CreateCharacterRequest)) {
    throw new Error('Expected argument of type fm.internal.CreateCharacterRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CreateCharacterRequest(buffer_arg) {
  return fminternal_internal_service_pb.CreateCharacterRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CreatePartyReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.CreatePartyReply)) {
    throw new Error('Expected argument of type fm.internal.CreatePartyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CreatePartyReply(buffer_arg) {
  return fminternal_internal_service_pb.CreatePartyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_CreatePartyRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.CreatePartyRequest)) {
    throw new Error('Expected argument of type fm.internal.CreatePartyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_CreatePartyRequest(buffer_arg) {
  return fminternal_internal_service_pb.CreatePartyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DeleteCharacterReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.DeleteCharacterReply)) {
    throw new Error('Expected argument of type fm.internal.DeleteCharacterReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DeleteCharacterReply(buffer_arg) {
  return fminternal_internal_service_pb.DeleteCharacterReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DeleteCharacterRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.DeleteCharacterRequest)) {
    throw new Error('Expected argument of type fm.internal.DeleteCharacterRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DeleteCharacterRequest(buffer_arg) {
  return fminternal_internal_service_pb.DeleteCharacterRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DenyPartyReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.DenyPartyReply)) {
    throw new Error('Expected argument of type fm.internal.DenyPartyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DenyPartyReply(buffer_arg) {
  return fminternal_internal_service_pb.DenyPartyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_DenyPartyRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.DenyPartyRequest)) {
    throw new Error('Expected argument of type fm.internal.DenyPartyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_DenyPartyRequest(buffer_arg) {
  return fminternal_internal_service_pb.DenyPartyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_EnterGameReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.EnterGameReply)) {
    throw new Error('Expected argument of type fm.internal.EnterGameReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_EnterGameReply(buffer_arg) {
  return fminternal_internal_service_pb.EnterGameReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_EnterGameRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.EnterGameRequest)) {
    throw new Error('Expected argument of type fm.internal.EnterGameRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_EnterGameRequest(buffer_arg) {
  return fminternal_internal_service_pb.EnterGameRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_ExpelPartyReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.ExpelPartyReply)) {
    throw new Error('Expected argument of type fm.internal.ExpelPartyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_ExpelPartyReply(buffer_arg) {
  return fminternal_internal_service_pb.ExpelPartyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_ExpelPartyRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.ExpelPartyRequest)) {
    throw new Error('Expected argument of type fm.internal.ExpelPartyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_ExpelPartyRequest(buffer_arg) {
  return fminternal_internal_service_pb.ExpelPartyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetCharacterListReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.GetCharacterListReply)) {
    throw new Error('Expected argument of type fm.internal.GetCharacterListReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetCharacterListReply(buffer_arg) {
  return fminternal_internal_service_pb.GetCharacterListReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetCharacterListRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.GetCharacterListRequest)) {
    throw new Error('Expected argument of type fm.internal.GetCharacterListRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetCharacterListRequest(buffer_arg) {
  return fminternal_internal_service_pb.GetCharacterListRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetPartyReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.GetPartyReply)) {
    throw new Error('Expected argument of type fm.internal.GetPartyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetPartyReply(buffer_arg) {
  return fminternal_internal_service_pb.GetPartyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetPartyRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.GetPartyRequest)) {
    throw new Error('Expected argument of type fm.internal.GetPartyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetPartyRequest(buffer_arg) {
  return fminternal_internal_service_pb.GetPartyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetServerCatalogReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.GetServerCatalogReply)) {
    throw new Error('Expected argument of type fm.internal.GetServerCatalogReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetServerCatalogReply(buffer_arg) {
  return fminternal_internal_service_pb.GetServerCatalogReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_GetServerCatalogRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.GetServerCatalogRequest)) {
    throw new Error('Expected argument of type fm.internal.GetServerCatalogRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_GetServerCatalogRequest(buffer_arg) {
  return fminternal_internal_service_pb.GetServerCatalogRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_InvitePartyReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.InvitePartyReply)) {
    throw new Error('Expected argument of type fm.internal.InvitePartyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_InvitePartyReply(buffer_arg) {
  return fminternal_internal_service_pb.InvitePartyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_InvitePartyRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.InvitePartyRequest)) {
    throw new Error('Expected argument of type fm.internal.InvitePartyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_InvitePartyRequest(buffer_arg) {
  return fminternal_internal_service_pb.InvitePartyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_JoinPartyReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.JoinPartyReply)) {
    throw new Error('Expected argument of type fm.internal.JoinPartyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_JoinPartyReply(buffer_arg) {
  return fminternal_internal_service_pb.JoinPartyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_JoinPartyRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.JoinPartyRequest)) {
    throw new Error('Expected argument of type fm.internal.JoinPartyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_JoinPartyRequest(buffer_arg) {
  return fminternal_internal_service_pb.JoinPartyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LeavePartyReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.LeavePartyReply)) {
    throw new Error('Expected argument of type fm.internal.LeavePartyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LeavePartyReply(buffer_arg) {
  return fminternal_internal_service_pb.LeavePartyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LeavePartyRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.LeavePartyRequest)) {
    throw new Error('Expected argument of type fm.internal.LeavePartyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LeavePartyRequest(buffer_arg) {
  return fminternal_internal_service_pb.LeavePartyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LoginAccountReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.LoginAccountReply)) {
    throw new Error('Expected argument of type fm.internal.LoginAccountReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LoginAccountReply(buffer_arg) {
  return fminternal_internal_service_pb.LoginAccountReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LoginAccountRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.LoginAccountRequest)) {
    throw new Error('Expected argument of type fm.internal.LoginAccountRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LoginAccountRequest(buffer_arg) {
  return fminternal_internal_service_pb.LoginAccountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LogoutSessionReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.LogoutSessionReply)) {
    throw new Error('Expected argument of type fm.internal.LogoutSessionReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LogoutSessionReply(buffer_arg) {
  return fminternal_internal_service_pb.LogoutSessionReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_LogoutSessionRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.LogoutSessionRequest)) {
    throw new Error('Expected argument of type fm.internal.LogoutSessionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_LogoutSessionRequest(buffer_arg) {
  return fminternal_internal_service_pb.LogoutSessionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_PingReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.PingReply)) {
    throw new Error('Expected argument of type fm.internal.PingReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_PingReply(buffer_arg) {
  return fminternal_internal_service_pb.PingReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_PingRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.PingRequest)) {
    throw new Error('Expected argument of type fm.internal.PingRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_PingRequest(buffer_arg) {
  return fminternal_internal_service_pb.PingRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_RefreshSessionReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.RefreshSessionReply)) {
    throw new Error('Expected argument of type fm.internal.RefreshSessionReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_RefreshSessionReply(buffer_arg) {
  return fminternal_internal_service_pb.RefreshSessionReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_RefreshSessionRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.RefreshSessionRequest)) {
    throw new Error('Expected argument of type fm.internal.RefreshSessionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_RefreshSessionRequest(buffer_arg) {
  return fminternal_internal_service_pb.RefreshSessionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveCharacterReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.SaveCharacterReply)) {
    throw new Error('Expected argument of type fm.internal.SaveCharacterReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharacterReply(buffer_arg) {
  return fminternal_internal_service_pb.SaveCharacterReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveCharacterRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.SaveCharacterRequest)) {
    throw new Error('Expected argument of type fm.internal.SaveCharacterRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharacterRequest(buffer_arg) {
  return fminternal_internal_service_pb.SaveCharacterRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveCharactersReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.SaveCharactersReply)) {
    throw new Error('Expected argument of type fm.internal.SaveCharactersReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharactersReply(buffer_arg) {
  return fminternal_internal_service_pb.SaveCharactersReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_SaveCharactersRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.SaveCharactersRequest)) {
    throw new Error('Expected argument of type fm.internal.SaveCharactersRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_SaveCharactersRequest(buffer_arg) {
  return fminternal_internal_service_pb.SaveCharactersRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_UpdatePartyMemberReply(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.UpdatePartyMemberReply)) {
    throw new Error('Expected argument of type fm.internal.UpdatePartyMemberReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_UpdatePartyMemberReply(buffer_arg) {
  return fminternal_internal_service_pb.UpdatePartyMemberReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_fm_internal_UpdatePartyMemberRequest(arg) {
  if (!(arg instanceof fminternal_internal_service_pb.UpdatePartyMemberRequest)) {
    throw new Error('Expected argument of type fm.internal.UpdatePartyMemberRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_fm_internal_UpdatePartyMemberRequest(buffer_arg) {
  return fminternal_internal_service_pb.UpdatePartyMemberRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var InternalService = exports.InternalService = {
  ping: {
    path: '/fm.internal.Internal/Ping',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.PingRequest,
    responseType: fminternal_internal_service_pb.PingReply,
    requestSerialize: serialize_fm_internal_PingRequest,
    requestDeserialize: deserialize_fm_internal_PingRequest,
    responseSerialize: serialize_fm_internal_PingReply,
    responseDeserialize: deserialize_fm_internal_PingReply,
  },
  getServerCatalog: {
    path: '/fm.internal.Internal/GetServerCatalog',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.GetServerCatalogRequest,
    responseType: fminternal_internal_service_pb.GetServerCatalogReply,
    requestSerialize: serialize_fm_internal_GetServerCatalogRequest,
    requestDeserialize: deserialize_fm_internal_GetServerCatalogRequest,
    responseSerialize: serialize_fm_internal_GetServerCatalogReply,
    responseDeserialize: deserialize_fm_internal_GetServerCatalogReply,
  },
  enterGame: {
    path: '/fm.internal.Internal/EnterGame',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.EnterGameRequest,
    responseType: fminternal_internal_service_pb.EnterGameReply,
    requestSerialize: serialize_fm_internal_EnterGameRequest,
    requestDeserialize: deserialize_fm_internal_EnterGameRequest,
    responseSerialize: serialize_fm_internal_EnterGameReply,
    responseDeserialize: deserialize_fm_internal_EnterGameReply,
  },
  beginGameTransition: {
    path: '/fm.internal.Internal/BeginGameTransition',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.BeginGameTransitionRequest,
    responseType: fminternal_internal_service_pb.BeginGameTransitionReply,
    requestSerialize: serialize_fm_internal_BeginGameTransitionRequest,
    requestDeserialize: deserialize_fm_internal_BeginGameTransitionRequest,
    responseSerialize: serialize_fm_internal_BeginGameTransitionReply,
    responseDeserialize: deserialize_fm_internal_BeginGameTransitionReply,
  },
  saveCharacter: {
    path: '/fm.internal.Internal/SaveCharacter',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.SaveCharacterRequest,
    responseType: fminternal_internal_service_pb.SaveCharacterReply,
    requestSerialize: serialize_fm_internal_SaveCharacterRequest,
    requestDeserialize: deserialize_fm_internal_SaveCharacterRequest,
    responseSerialize: serialize_fm_internal_SaveCharacterReply,
    responseDeserialize: deserialize_fm_internal_SaveCharacterReply,
  },
  saveCharacters: {
    path: '/fm.internal.Internal/SaveCharacters',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.SaveCharactersRequest,
    responseType: fminternal_internal_service_pb.SaveCharactersReply,
    requestSerialize: serialize_fm_internal_SaveCharactersRequest,
    requestDeserialize: deserialize_fm_internal_SaveCharactersRequest,
    responseSerialize: serialize_fm_internal_SaveCharactersReply,
    responseDeserialize: deserialize_fm_internal_SaveCharactersReply,
  },
  loginAccount: {
    path: '/fm.internal.Internal/LoginAccount',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.LoginAccountRequest,
    responseType: fminternal_internal_service_pb.LoginAccountReply,
    requestSerialize: serialize_fm_internal_LoginAccountRequest,
    requestDeserialize: deserialize_fm_internal_LoginAccountRequest,
    responseSerialize: serialize_fm_internal_LoginAccountReply,
    responseDeserialize: deserialize_fm_internal_LoginAccountReply,
  },
  getCharacterList: {
    path: '/fm.internal.Internal/GetCharacterList',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.GetCharacterListRequest,
    responseType: fminternal_internal_service_pb.GetCharacterListReply,
    requestSerialize: serialize_fm_internal_GetCharacterListRequest,
    requestDeserialize: deserialize_fm_internal_GetCharacterListRequest,
    responseSerialize: serialize_fm_internal_GetCharacterListReply,
    responseDeserialize: deserialize_fm_internal_GetCharacterListReply,
  },
  checkCharacterName: {
    path: '/fm.internal.Internal/CheckCharacterName',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.CheckCharacterNameRequest,
    responseType: fminternal_internal_service_pb.CheckCharacterNameReply,
    requestSerialize: serialize_fm_internal_CheckCharacterNameRequest,
    requestDeserialize: deserialize_fm_internal_CheckCharacterNameRequest,
    responseSerialize: serialize_fm_internal_CheckCharacterNameReply,
    responseDeserialize: deserialize_fm_internal_CheckCharacterNameReply,
  },
  createCharacter: {
    path: '/fm.internal.Internal/CreateCharacter',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.CreateCharacterRequest,
    responseType: fminternal_internal_service_pb.CreateCharacterReply,
    requestSerialize: serialize_fm_internal_CreateCharacterRequest,
    requestDeserialize: deserialize_fm_internal_CreateCharacterRequest,
    responseSerialize: serialize_fm_internal_CreateCharacterReply,
    responseDeserialize: deserialize_fm_internal_CreateCharacterReply,
  },
  deleteCharacter: {
    path: '/fm.internal.Internal/DeleteCharacter',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.DeleteCharacterRequest,
    responseType: fminternal_internal_service_pb.DeleteCharacterReply,
    requestSerialize: serialize_fm_internal_DeleteCharacterRequest,
    requestDeserialize: deserialize_fm_internal_DeleteCharacterRequest,
    responseSerialize: serialize_fm_internal_DeleteCharacterReply,
    responseDeserialize: deserialize_fm_internal_DeleteCharacterReply,
  },
  refreshSession: {
    path: '/fm.internal.Internal/RefreshSession',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.RefreshSessionRequest,
    responseType: fminternal_internal_service_pb.RefreshSessionReply,
    requestSerialize: serialize_fm_internal_RefreshSessionRequest,
    requestDeserialize: deserialize_fm_internal_RefreshSessionRequest,
    responseSerialize: serialize_fm_internal_RefreshSessionReply,
    responseDeserialize: deserialize_fm_internal_RefreshSessionReply,
  },
  logoutSession: {
    path: '/fm.internal.Internal/LogoutSession',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.LogoutSessionRequest,
    responseType: fminternal_internal_service_pb.LogoutSessionReply,
    requestSerialize: serialize_fm_internal_LogoutSessionRequest,
    requestDeserialize: deserialize_fm_internal_LogoutSessionRequest,
    responseSerialize: serialize_fm_internal_LogoutSessionReply,
    responseDeserialize: deserialize_fm_internal_LogoutSessionReply,
  },
  createParty: {
    path: '/fm.internal.Internal/CreateParty',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.CreatePartyRequest,
    responseType: fminternal_internal_service_pb.CreatePartyReply,
    requestSerialize: serialize_fm_internal_CreatePartyRequest,
    requestDeserialize: deserialize_fm_internal_CreatePartyRequest,
    responseSerialize: serialize_fm_internal_CreatePartyReply,
    responseDeserialize: deserialize_fm_internal_CreatePartyReply,
  },
  joinParty: {
    path: '/fm.internal.Internal/JoinParty',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.JoinPartyRequest,
    responseType: fminternal_internal_service_pb.JoinPartyReply,
    requestSerialize: serialize_fm_internal_JoinPartyRequest,
    requestDeserialize: deserialize_fm_internal_JoinPartyRequest,
    responseSerialize: serialize_fm_internal_JoinPartyReply,
    responseDeserialize: deserialize_fm_internal_JoinPartyReply,
  },
  leaveParty: {
    path: '/fm.internal.Internal/LeaveParty',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.LeavePartyRequest,
    responseType: fminternal_internal_service_pb.LeavePartyReply,
    requestSerialize: serialize_fm_internal_LeavePartyRequest,
    requestDeserialize: deserialize_fm_internal_LeavePartyRequest,
    responseSerialize: serialize_fm_internal_LeavePartyReply,
    responseDeserialize: deserialize_fm_internal_LeavePartyReply,
  },
  expelParty: {
    path: '/fm.internal.Internal/ExpelParty',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.ExpelPartyRequest,
    responseType: fminternal_internal_service_pb.ExpelPartyReply,
    requestSerialize: serialize_fm_internal_ExpelPartyRequest,
    requestDeserialize: deserialize_fm_internal_ExpelPartyRequest,
    responseSerialize: serialize_fm_internal_ExpelPartyReply,
    responseDeserialize: deserialize_fm_internal_ExpelPartyReply,
  },
  changePartyLeader: {
    path: '/fm.internal.Internal/ChangePartyLeader',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.ChangePartyLeaderRequest,
    responseType: fminternal_internal_service_pb.ChangePartyLeaderReply,
    requestSerialize: serialize_fm_internal_ChangePartyLeaderRequest,
    requestDeserialize: deserialize_fm_internal_ChangePartyLeaderRequest,
    responseSerialize: serialize_fm_internal_ChangePartyLeaderReply,
    responseDeserialize: deserialize_fm_internal_ChangePartyLeaderReply,
  },
  getParty: {
    path: '/fm.internal.Internal/GetParty',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.GetPartyRequest,
    responseType: fminternal_internal_service_pb.GetPartyReply,
    requestSerialize: serialize_fm_internal_GetPartyRequest,
    requestDeserialize: deserialize_fm_internal_GetPartyRequest,
    responseSerialize: serialize_fm_internal_GetPartyReply,
    responseDeserialize: deserialize_fm_internal_GetPartyReply,
  },
  updatePartyMember: {
    path: '/fm.internal.Internal/UpdatePartyMember',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.UpdatePartyMemberRequest,
    responseType: fminternal_internal_service_pb.UpdatePartyMemberReply,
    requestSerialize: serialize_fm_internal_UpdatePartyMemberRequest,
    requestDeserialize: deserialize_fm_internal_UpdatePartyMemberRequest,
    responseSerialize: serialize_fm_internal_UpdatePartyMemberReply,
    responseDeserialize: deserialize_fm_internal_UpdatePartyMemberReply,
  },
  inviteParty: {
    path: '/fm.internal.Internal/InviteParty',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.InvitePartyRequest,
    responseType: fminternal_internal_service_pb.InvitePartyReply,
    requestSerialize: serialize_fm_internal_InvitePartyRequest,
    requestDeserialize: deserialize_fm_internal_InvitePartyRequest,
    responseSerialize: serialize_fm_internal_InvitePartyReply,
    responseDeserialize: deserialize_fm_internal_InvitePartyReply,
  },
  denyParty: {
    path: '/fm.internal.Internal/DenyParty',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.DenyPartyRequest,
    responseType: fminternal_internal_service_pb.DenyPartyReply,
    requestSerialize: serialize_fm_internal_DenyPartyRequest,
    requestDeserialize: deserialize_fm_internal_DenyPartyRequest,
    responseSerialize: serialize_fm_internal_DenyPartyReply,
    responseDeserialize: deserialize_fm_internal_DenyPartyReply,
  },
  broadcastPartyChat: {
    path: '/fm.internal.Internal/BroadcastPartyChat',
    requestStream: false,
    responseStream: false,
    requestType: fminternal_internal_service_pb.BroadcastPartyChatRequest,
    responseType: fminternal_internal_service_pb.BroadcastPartyChatReply,
    requestSerialize: serialize_fm_internal_BroadcastPartyChatRequest,
    requestDeserialize: deserialize_fm_internal_BroadcastPartyChatRequest,
    responseSerialize: serialize_fm_internal_BroadcastPartyChatReply,
    responseDeserialize: deserialize_fm_internal_BroadcastPartyChatReply,
  },
};

exports.InternalClient = grpc.makeGenericClientConstructor(InternalService, 'Internal');
