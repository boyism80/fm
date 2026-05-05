// package: fm.internal
// file: fminternal/internal_service.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as fminternal_internal_service_pb from "../fminternal/internal_service_pb";

interface IInternalService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    ping: IInternalService_IPing;
    getServerCatalog: IInternalService_IGetServerCatalog;
    enterGame: IInternalService_IEnterGame;
    beginGameTransition: IInternalService_IBeginGameTransition;
    saveCharacter: IInternalService_ISaveCharacter;
    saveCharacters: IInternalService_ISaveCharacters;
    loginAccount: IInternalService_ILoginAccount;
    getCharacterList: IInternalService_IGetCharacterList;
    checkCharacterName: IInternalService_ICheckCharacterName;
    createCharacter: IInternalService_ICreateCharacter;
    deleteCharacter: IInternalService_IDeleteCharacter;
    refreshSession: IInternalService_IRefreshSession;
    logoutSession: IInternalService_ILogoutSession;
    createParty: IInternalService_ICreateParty;
    joinParty: IInternalService_IJoinParty;
    leaveParty: IInternalService_ILeaveParty;
    expelParty: IInternalService_IExpelParty;
    changePartyLeader: IInternalService_IChangePartyLeader;
    getParty: IInternalService_IGetParty;
    updatePartyMember: IInternalService_IUpdatePartyMember;
    inviteParty: IInternalService_IInviteParty;
    denyParty: IInternalService_IDenyParty;
    broadcastMultiChat: IInternalService_IBroadcastMultiChat;
}

interface IInternalService_IPing extends grpc.MethodDefinition<fminternal_internal_service_pb.PingRequest, fminternal_internal_service_pb.PingReply> {
    path: "/fm.internal.Internal/Ping";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.PingRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.PingRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.PingReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.PingReply>;
}
interface IInternalService_IGetServerCatalog extends grpc.MethodDefinition<fminternal_internal_service_pb.GetServerCatalogRequest, fminternal_internal_service_pb.GetServerCatalogReply> {
    path: "/fm.internal.Internal/GetServerCatalog";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.GetServerCatalogRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.GetServerCatalogRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.GetServerCatalogReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.GetServerCatalogReply>;
}
interface IInternalService_IEnterGame extends grpc.MethodDefinition<fminternal_internal_service_pb.EnterGameRequest, fminternal_internal_service_pb.EnterGameReply> {
    path: "/fm.internal.Internal/EnterGame";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.EnterGameRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.EnterGameRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.EnterGameReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.EnterGameReply>;
}
interface IInternalService_IBeginGameTransition extends grpc.MethodDefinition<fminternal_internal_service_pb.BeginGameTransitionRequest, fminternal_internal_service_pb.BeginGameTransitionReply> {
    path: "/fm.internal.Internal/BeginGameTransition";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.BeginGameTransitionRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.BeginGameTransitionRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.BeginGameTransitionReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.BeginGameTransitionReply>;
}
interface IInternalService_ISaveCharacter extends grpc.MethodDefinition<fminternal_internal_service_pb.SaveCharacterRequest, fminternal_internal_service_pb.SaveCharacterReply> {
    path: "/fm.internal.Internal/SaveCharacter";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.SaveCharacterRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.SaveCharacterRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.SaveCharacterReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.SaveCharacterReply>;
}
interface IInternalService_ISaveCharacters extends grpc.MethodDefinition<fminternal_internal_service_pb.SaveCharactersRequest, fminternal_internal_service_pb.SaveCharactersReply> {
    path: "/fm.internal.Internal/SaveCharacters";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.SaveCharactersRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.SaveCharactersRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.SaveCharactersReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.SaveCharactersReply>;
}
interface IInternalService_ILoginAccount extends grpc.MethodDefinition<fminternal_internal_service_pb.LoginAccountRequest, fminternal_internal_service_pb.LoginAccountReply> {
    path: "/fm.internal.Internal/LoginAccount";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.LoginAccountRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.LoginAccountRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.LoginAccountReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.LoginAccountReply>;
}
interface IInternalService_IGetCharacterList extends grpc.MethodDefinition<fminternal_internal_service_pb.GetCharacterListRequest, fminternal_internal_service_pb.GetCharacterListReply> {
    path: "/fm.internal.Internal/GetCharacterList";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.GetCharacterListRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.GetCharacterListRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.GetCharacterListReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.GetCharacterListReply>;
}
interface IInternalService_ICheckCharacterName extends grpc.MethodDefinition<fminternal_internal_service_pb.CheckCharacterNameRequest, fminternal_internal_service_pb.CheckCharacterNameReply> {
    path: "/fm.internal.Internal/CheckCharacterName";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.CheckCharacterNameRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.CheckCharacterNameRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.CheckCharacterNameReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.CheckCharacterNameReply>;
}
interface IInternalService_ICreateCharacter extends grpc.MethodDefinition<fminternal_internal_service_pb.CreateCharacterRequest, fminternal_internal_service_pb.CreateCharacterReply> {
    path: "/fm.internal.Internal/CreateCharacter";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.CreateCharacterRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.CreateCharacterRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.CreateCharacterReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.CreateCharacterReply>;
}
interface IInternalService_IDeleteCharacter extends grpc.MethodDefinition<fminternal_internal_service_pb.DeleteCharacterRequest, fminternal_internal_service_pb.DeleteCharacterReply> {
    path: "/fm.internal.Internal/DeleteCharacter";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.DeleteCharacterRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.DeleteCharacterRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.DeleteCharacterReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.DeleteCharacterReply>;
}
interface IInternalService_IRefreshSession extends grpc.MethodDefinition<fminternal_internal_service_pb.RefreshSessionRequest, fminternal_internal_service_pb.RefreshSessionReply> {
    path: "/fm.internal.Internal/RefreshSession";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.RefreshSessionRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.RefreshSessionRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.RefreshSessionReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.RefreshSessionReply>;
}
interface IInternalService_ILogoutSession extends grpc.MethodDefinition<fminternal_internal_service_pb.LogoutSessionRequest, fminternal_internal_service_pb.LogoutSessionReply> {
    path: "/fm.internal.Internal/LogoutSession";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.LogoutSessionRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.LogoutSessionRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.LogoutSessionReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.LogoutSessionReply>;
}
interface IInternalService_ICreateParty extends grpc.MethodDefinition<fminternal_internal_service_pb.CreatePartyRequest, fminternal_internal_service_pb.CreatePartyReply> {
    path: "/fm.internal.Internal/CreateParty";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.CreatePartyRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.CreatePartyRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.CreatePartyReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.CreatePartyReply>;
}
interface IInternalService_IJoinParty extends grpc.MethodDefinition<fminternal_internal_service_pb.JoinPartyRequest, fminternal_internal_service_pb.JoinPartyReply> {
    path: "/fm.internal.Internal/JoinParty";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.JoinPartyRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.JoinPartyRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.JoinPartyReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.JoinPartyReply>;
}
interface IInternalService_ILeaveParty extends grpc.MethodDefinition<fminternal_internal_service_pb.LeavePartyRequest, fminternal_internal_service_pb.LeavePartyReply> {
    path: "/fm.internal.Internal/LeaveParty";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.LeavePartyRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.LeavePartyRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.LeavePartyReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.LeavePartyReply>;
}
interface IInternalService_IExpelParty extends grpc.MethodDefinition<fminternal_internal_service_pb.ExpelPartyRequest, fminternal_internal_service_pb.ExpelPartyReply> {
    path: "/fm.internal.Internal/ExpelParty";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.ExpelPartyRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.ExpelPartyRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.ExpelPartyReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.ExpelPartyReply>;
}
interface IInternalService_IChangePartyLeader extends grpc.MethodDefinition<fminternal_internal_service_pb.ChangePartyLeaderRequest, fminternal_internal_service_pb.ChangePartyLeaderReply> {
    path: "/fm.internal.Internal/ChangePartyLeader";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.ChangePartyLeaderRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.ChangePartyLeaderRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.ChangePartyLeaderReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.ChangePartyLeaderReply>;
}
interface IInternalService_IGetParty extends grpc.MethodDefinition<fminternal_internal_service_pb.GetPartyRequest, fminternal_internal_service_pb.GetPartyReply> {
    path: "/fm.internal.Internal/GetParty";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.GetPartyRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.GetPartyRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.GetPartyReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.GetPartyReply>;
}
interface IInternalService_IUpdatePartyMember extends grpc.MethodDefinition<fminternal_internal_service_pb.UpdatePartyMemberRequest, fminternal_internal_service_pb.UpdatePartyMemberReply> {
    path: "/fm.internal.Internal/UpdatePartyMember";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.UpdatePartyMemberRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.UpdatePartyMemberRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.UpdatePartyMemberReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.UpdatePartyMemberReply>;
}
interface IInternalService_IInviteParty extends grpc.MethodDefinition<fminternal_internal_service_pb.InvitePartyRequest, fminternal_internal_service_pb.InvitePartyReply> {
    path: "/fm.internal.Internal/InviteParty";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.InvitePartyRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.InvitePartyRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.InvitePartyReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.InvitePartyReply>;
}
interface IInternalService_IDenyParty extends grpc.MethodDefinition<fminternal_internal_service_pb.DenyPartyRequest, fminternal_internal_service_pb.DenyPartyReply> {
    path: "/fm.internal.Internal/DenyParty";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.DenyPartyRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.DenyPartyRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.DenyPartyReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.DenyPartyReply>;
}
interface IInternalService_IBroadcastMultiChat extends grpc.MethodDefinition<fminternal_internal_service_pb.BroadcastMultiChatRequest, fminternal_internal_service_pb.BroadcastMultiChatReply> {
    path: "/fm.internal.Internal/BroadcastMultiChat";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<fminternal_internal_service_pb.BroadcastMultiChatRequest>;
    requestDeserialize: grpc.deserialize<fminternal_internal_service_pb.BroadcastMultiChatRequest>;
    responseSerialize: grpc.serialize<fminternal_internal_service_pb.BroadcastMultiChatReply>;
    responseDeserialize: grpc.deserialize<fminternal_internal_service_pb.BroadcastMultiChatReply>;
}

export const InternalService: IInternalService;

export interface IInternalServer extends grpc.UntypedServiceImplementation {
    ping: grpc.handleUnaryCall<fminternal_internal_service_pb.PingRequest, fminternal_internal_service_pb.PingReply>;
    getServerCatalog: grpc.handleUnaryCall<fminternal_internal_service_pb.GetServerCatalogRequest, fminternal_internal_service_pb.GetServerCatalogReply>;
    enterGame: grpc.handleUnaryCall<fminternal_internal_service_pb.EnterGameRequest, fminternal_internal_service_pb.EnterGameReply>;
    beginGameTransition: grpc.handleUnaryCall<fminternal_internal_service_pb.BeginGameTransitionRequest, fminternal_internal_service_pb.BeginGameTransitionReply>;
    saveCharacter: grpc.handleUnaryCall<fminternal_internal_service_pb.SaveCharacterRequest, fminternal_internal_service_pb.SaveCharacterReply>;
    saveCharacters: grpc.handleUnaryCall<fminternal_internal_service_pb.SaveCharactersRequest, fminternal_internal_service_pb.SaveCharactersReply>;
    loginAccount: grpc.handleUnaryCall<fminternal_internal_service_pb.LoginAccountRequest, fminternal_internal_service_pb.LoginAccountReply>;
    getCharacterList: grpc.handleUnaryCall<fminternal_internal_service_pb.GetCharacterListRequest, fminternal_internal_service_pb.GetCharacterListReply>;
    checkCharacterName: grpc.handleUnaryCall<fminternal_internal_service_pb.CheckCharacterNameRequest, fminternal_internal_service_pb.CheckCharacterNameReply>;
    createCharacter: grpc.handleUnaryCall<fminternal_internal_service_pb.CreateCharacterRequest, fminternal_internal_service_pb.CreateCharacterReply>;
    deleteCharacter: grpc.handleUnaryCall<fminternal_internal_service_pb.DeleteCharacterRequest, fminternal_internal_service_pb.DeleteCharacterReply>;
    refreshSession: grpc.handleUnaryCall<fminternal_internal_service_pb.RefreshSessionRequest, fminternal_internal_service_pb.RefreshSessionReply>;
    logoutSession: grpc.handleUnaryCall<fminternal_internal_service_pb.LogoutSessionRequest, fminternal_internal_service_pb.LogoutSessionReply>;
    createParty: grpc.handleUnaryCall<fminternal_internal_service_pb.CreatePartyRequest, fminternal_internal_service_pb.CreatePartyReply>;
    joinParty: grpc.handleUnaryCall<fminternal_internal_service_pb.JoinPartyRequest, fminternal_internal_service_pb.JoinPartyReply>;
    leaveParty: grpc.handleUnaryCall<fminternal_internal_service_pb.LeavePartyRequest, fminternal_internal_service_pb.LeavePartyReply>;
    expelParty: grpc.handleUnaryCall<fminternal_internal_service_pb.ExpelPartyRequest, fminternal_internal_service_pb.ExpelPartyReply>;
    changePartyLeader: grpc.handleUnaryCall<fminternal_internal_service_pb.ChangePartyLeaderRequest, fminternal_internal_service_pb.ChangePartyLeaderReply>;
    getParty: grpc.handleUnaryCall<fminternal_internal_service_pb.GetPartyRequest, fminternal_internal_service_pb.GetPartyReply>;
    updatePartyMember: grpc.handleUnaryCall<fminternal_internal_service_pb.UpdatePartyMemberRequest, fminternal_internal_service_pb.UpdatePartyMemberReply>;
    inviteParty: grpc.handleUnaryCall<fminternal_internal_service_pb.InvitePartyRequest, fminternal_internal_service_pb.InvitePartyReply>;
    denyParty: grpc.handleUnaryCall<fminternal_internal_service_pb.DenyPartyRequest, fminternal_internal_service_pb.DenyPartyReply>;
    broadcastMultiChat: grpc.handleUnaryCall<fminternal_internal_service_pb.BroadcastMultiChatRequest, fminternal_internal_service_pb.BroadcastMultiChatReply>;
}

export interface IInternalClient {
    ping(request: fminternal_internal_service_pb.PingRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.PingReply) => void): grpc.ClientUnaryCall;
    ping(request: fminternal_internal_service_pb.PingRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.PingReply) => void): grpc.ClientUnaryCall;
    ping(request: fminternal_internal_service_pb.PingRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.PingReply) => void): grpc.ClientUnaryCall;
    getServerCatalog(request: fminternal_internal_service_pb.GetServerCatalogRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetServerCatalogReply) => void): grpc.ClientUnaryCall;
    getServerCatalog(request: fminternal_internal_service_pb.GetServerCatalogRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetServerCatalogReply) => void): grpc.ClientUnaryCall;
    getServerCatalog(request: fminternal_internal_service_pb.GetServerCatalogRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetServerCatalogReply) => void): grpc.ClientUnaryCall;
    enterGame(request: fminternal_internal_service_pb.EnterGameRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.EnterGameReply) => void): grpc.ClientUnaryCall;
    enterGame(request: fminternal_internal_service_pb.EnterGameRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.EnterGameReply) => void): grpc.ClientUnaryCall;
    enterGame(request: fminternal_internal_service_pb.EnterGameRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.EnterGameReply) => void): grpc.ClientUnaryCall;
    beginGameTransition(request: fminternal_internal_service_pb.BeginGameTransitionRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BeginGameTransitionReply) => void): grpc.ClientUnaryCall;
    beginGameTransition(request: fminternal_internal_service_pb.BeginGameTransitionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BeginGameTransitionReply) => void): grpc.ClientUnaryCall;
    beginGameTransition(request: fminternal_internal_service_pb.BeginGameTransitionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BeginGameTransitionReply) => void): grpc.ClientUnaryCall;
    saveCharacter(request: fminternal_internal_service_pb.SaveCharacterRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharacterReply) => void): grpc.ClientUnaryCall;
    saveCharacter(request: fminternal_internal_service_pb.SaveCharacterRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharacterReply) => void): grpc.ClientUnaryCall;
    saveCharacter(request: fminternal_internal_service_pb.SaveCharacterRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharacterReply) => void): grpc.ClientUnaryCall;
    saveCharacters(request: fminternal_internal_service_pb.SaveCharactersRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharactersReply) => void): grpc.ClientUnaryCall;
    saveCharacters(request: fminternal_internal_service_pb.SaveCharactersRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharactersReply) => void): grpc.ClientUnaryCall;
    saveCharacters(request: fminternal_internal_service_pb.SaveCharactersRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharactersReply) => void): grpc.ClientUnaryCall;
    loginAccount(request: fminternal_internal_service_pb.LoginAccountRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LoginAccountReply) => void): grpc.ClientUnaryCall;
    loginAccount(request: fminternal_internal_service_pb.LoginAccountRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LoginAccountReply) => void): grpc.ClientUnaryCall;
    loginAccount(request: fminternal_internal_service_pb.LoginAccountRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LoginAccountReply) => void): grpc.ClientUnaryCall;
    getCharacterList(request: fminternal_internal_service_pb.GetCharacterListRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetCharacterListReply) => void): grpc.ClientUnaryCall;
    getCharacterList(request: fminternal_internal_service_pb.GetCharacterListRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetCharacterListReply) => void): grpc.ClientUnaryCall;
    getCharacterList(request: fminternal_internal_service_pb.GetCharacterListRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetCharacterListReply) => void): grpc.ClientUnaryCall;
    checkCharacterName(request: fminternal_internal_service_pb.CheckCharacterNameRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CheckCharacterNameReply) => void): grpc.ClientUnaryCall;
    checkCharacterName(request: fminternal_internal_service_pb.CheckCharacterNameRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CheckCharacterNameReply) => void): grpc.ClientUnaryCall;
    checkCharacterName(request: fminternal_internal_service_pb.CheckCharacterNameRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CheckCharacterNameReply) => void): grpc.ClientUnaryCall;
    createCharacter(request: fminternal_internal_service_pb.CreateCharacterRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreateCharacterReply) => void): grpc.ClientUnaryCall;
    createCharacter(request: fminternal_internal_service_pb.CreateCharacterRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreateCharacterReply) => void): grpc.ClientUnaryCall;
    createCharacter(request: fminternal_internal_service_pb.CreateCharacterRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreateCharacterReply) => void): grpc.ClientUnaryCall;
    deleteCharacter(request: fminternal_internal_service_pb.DeleteCharacterRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DeleteCharacterReply) => void): grpc.ClientUnaryCall;
    deleteCharacter(request: fminternal_internal_service_pb.DeleteCharacterRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DeleteCharacterReply) => void): grpc.ClientUnaryCall;
    deleteCharacter(request: fminternal_internal_service_pb.DeleteCharacterRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DeleteCharacterReply) => void): grpc.ClientUnaryCall;
    refreshSession(request: fminternal_internal_service_pb.RefreshSessionRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.RefreshSessionReply) => void): grpc.ClientUnaryCall;
    refreshSession(request: fminternal_internal_service_pb.RefreshSessionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.RefreshSessionReply) => void): grpc.ClientUnaryCall;
    refreshSession(request: fminternal_internal_service_pb.RefreshSessionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.RefreshSessionReply) => void): grpc.ClientUnaryCall;
    logoutSession(request: fminternal_internal_service_pb.LogoutSessionRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LogoutSessionReply) => void): grpc.ClientUnaryCall;
    logoutSession(request: fminternal_internal_service_pb.LogoutSessionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LogoutSessionReply) => void): grpc.ClientUnaryCall;
    logoutSession(request: fminternal_internal_service_pb.LogoutSessionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LogoutSessionReply) => void): grpc.ClientUnaryCall;
    createParty(request: fminternal_internal_service_pb.CreatePartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreatePartyReply) => void): grpc.ClientUnaryCall;
    createParty(request: fminternal_internal_service_pb.CreatePartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreatePartyReply) => void): grpc.ClientUnaryCall;
    createParty(request: fminternal_internal_service_pb.CreatePartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreatePartyReply) => void): grpc.ClientUnaryCall;
    joinParty(request: fminternal_internal_service_pb.JoinPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.JoinPartyReply) => void): grpc.ClientUnaryCall;
    joinParty(request: fminternal_internal_service_pb.JoinPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.JoinPartyReply) => void): grpc.ClientUnaryCall;
    joinParty(request: fminternal_internal_service_pb.JoinPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.JoinPartyReply) => void): grpc.ClientUnaryCall;
    leaveParty(request: fminternal_internal_service_pb.LeavePartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LeavePartyReply) => void): grpc.ClientUnaryCall;
    leaveParty(request: fminternal_internal_service_pb.LeavePartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LeavePartyReply) => void): grpc.ClientUnaryCall;
    leaveParty(request: fminternal_internal_service_pb.LeavePartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LeavePartyReply) => void): grpc.ClientUnaryCall;
    expelParty(request: fminternal_internal_service_pb.ExpelPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ExpelPartyReply) => void): grpc.ClientUnaryCall;
    expelParty(request: fminternal_internal_service_pb.ExpelPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ExpelPartyReply) => void): grpc.ClientUnaryCall;
    expelParty(request: fminternal_internal_service_pb.ExpelPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ExpelPartyReply) => void): grpc.ClientUnaryCall;
    changePartyLeader(request: fminternal_internal_service_pb.ChangePartyLeaderRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ChangePartyLeaderReply) => void): grpc.ClientUnaryCall;
    changePartyLeader(request: fminternal_internal_service_pb.ChangePartyLeaderRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ChangePartyLeaderReply) => void): grpc.ClientUnaryCall;
    changePartyLeader(request: fminternal_internal_service_pb.ChangePartyLeaderRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ChangePartyLeaderReply) => void): grpc.ClientUnaryCall;
    getParty(request: fminternal_internal_service_pb.GetPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetPartyReply) => void): grpc.ClientUnaryCall;
    getParty(request: fminternal_internal_service_pb.GetPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetPartyReply) => void): grpc.ClientUnaryCall;
    getParty(request: fminternal_internal_service_pb.GetPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetPartyReply) => void): grpc.ClientUnaryCall;
    updatePartyMember(request: fminternal_internal_service_pb.UpdatePartyMemberRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.UpdatePartyMemberReply) => void): grpc.ClientUnaryCall;
    updatePartyMember(request: fminternal_internal_service_pb.UpdatePartyMemberRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.UpdatePartyMemberReply) => void): grpc.ClientUnaryCall;
    updatePartyMember(request: fminternal_internal_service_pb.UpdatePartyMemberRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.UpdatePartyMemberReply) => void): grpc.ClientUnaryCall;
    inviteParty(request: fminternal_internal_service_pb.InvitePartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.InvitePartyReply) => void): grpc.ClientUnaryCall;
    inviteParty(request: fminternal_internal_service_pb.InvitePartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.InvitePartyReply) => void): grpc.ClientUnaryCall;
    inviteParty(request: fminternal_internal_service_pb.InvitePartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.InvitePartyReply) => void): grpc.ClientUnaryCall;
    denyParty(request: fminternal_internal_service_pb.DenyPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DenyPartyReply) => void): grpc.ClientUnaryCall;
    denyParty(request: fminternal_internal_service_pb.DenyPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DenyPartyReply) => void): grpc.ClientUnaryCall;
    denyParty(request: fminternal_internal_service_pb.DenyPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DenyPartyReply) => void): grpc.ClientUnaryCall;
    broadcastMultiChat(request: fminternal_internal_service_pb.BroadcastMultiChatRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BroadcastMultiChatReply) => void): grpc.ClientUnaryCall;
    broadcastMultiChat(request: fminternal_internal_service_pb.BroadcastMultiChatRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BroadcastMultiChatReply) => void): grpc.ClientUnaryCall;
    broadcastMultiChat(request: fminternal_internal_service_pb.BroadcastMultiChatRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BroadcastMultiChatReply) => void): grpc.ClientUnaryCall;
}

export class InternalClient extends grpc.Client implements IInternalClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public ping(request: fminternal_internal_service_pb.PingRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.PingReply) => void): grpc.ClientUnaryCall;
    public ping(request: fminternal_internal_service_pb.PingRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.PingReply) => void): grpc.ClientUnaryCall;
    public ping(request: fminternal_internal_service_pb.PingRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.PingReply) => void): grpc.ClientUnaryCall;
    public getServerCatalog(request: fminternal_internal_service_pb.GetServerCatalogRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetServerCatalogReply) => void): grpc.ClientUnaryCall;
    public getServerCatalog(request: fminternal_internal_service_pb.GetServerCatalogRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetServerCatalogReply) => void): grpc.ClientUnaryCall;
    public getServerCatalog(request: fminternal_internal_service_pb.GetServerCatalogRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetServerCatalogReply) => void): grpc.ClientUnaryCall;
    public enterGame(request: fminternal_internal_service_pb.EnterGameRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.EnterGameReply) => void): grpc.ClientUnaryCall;
    public enterGame(request: fminternal_internal_service_pb.EnterGameRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.EnterGameReply) => void): grpc.ClientUnaryCall;
    public enterGame(request: fminternal_internal_service_pb.EnterGameRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.EnterGameReply) => void): grpc.ClientUnaryCall;
    public beginGameTransition(request: fminternal_internal_service_pb.BeginGameTransitionRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BeginGameTransitionReply) => void): grpc.ClientUnaryCall;
    public beginGameTransition(request: fminternal_internal_service_pb.BeginGameTransitionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BeginGameTransitionReply) => void): grpc.ClientUnaryCall;
    public beginGameTransition(request: fminternal_internal_service_pb.BeginGameTransitionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BeginGameTransitionReply) => void): grpc.ClientUnaryCall;
    public saveCharacter(request: fminternal_internal_service_pb.SaveCharacterRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharacterReply) => void): grpc.ClientUnaryCall;
    public saveCharacter(request: fminternal_internal_service_pb.SaveCharacterRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharacterReply) => void): grpc.ClientUnaryCall;
    public saveCharacter(request: fminternal_internal_service_pb.SaveCharacterRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharacterReply) => void): grpc.ClientUnaryCall;
    public saveCharacters(request: fminternal_internal_service_pb.SaveCharactersRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharactersReply) => void): grpc.ClientUnaryCall;
    public saveCharacters(request: fminternal_internal_service_pb.SaveCharactersRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharactersReply) => void): grpc.ClientUnaryCall;
    public saveCharacters(request: fminternal_internal_service_pb.SaveCharactersRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.SaveCharactersReply) => void): grpc.ClientUnaryCall;
    public loginAccount(request: fminternal_internal_service_pb.LoginAccountRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LoginAccountReply) => void): grpc.ClientUnaryCall;
    public loginAccount(request: fminternal_internal_service_pb.LoginAccountRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LoginAccountReply) => void): grpc.ClientUnaryCall;
    public loginAccount(request: fminternal_internal_service_pb.LoginAccountRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LoginAccountReply) => void): grpc.ClientUnaryCall;
    public getCharacterList(request: fminternal_internal_service_pb.GetCharacterListRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetCharacterListReply) => void): grpc.ClientUnaryCall;
    public getCharacterList(request: fminternal_internal_service_pb.GetCharacterListRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetCharacterListReply) => void): grpc.ClientUnaryCall;
    public getCharacterList(request: fminternal_internal_service_pb.GetCharacterListRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetCharacterListReply) => void): grpc.ClientUnaryCall;
    public checkCharacterName(request: fminternal_internal_service_pb.CheckCharacterNameRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CheckCharacterNameReply) => void): grpc.ClientUnaryCall;
    public checkCharacterName(request: fminternal_internal_service_pb.CheckCharacterNameRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CheckCharacterNameReply) => void): grpc.ClientUnaryCall;
    public checkCharacterName(request: fminternal_internal_service_pb.CheckCharacterNameRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CheckCharacterNameReply) => void): grpc.ClientUnaryCall;
    public createCharacter(request: fminternal_internal_service_pb.CreateCharacterRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreateCharacterReply) => void): grpc.ClientUnaryCall;
    public createCharacter(request: fminternal_internal_service_pb.CreateCharacterRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreateCharacterReply) => void): grpc.ClientUnaryCall;
    public createCharacter(request: fminternal_internal_service_pb.CreateCharacterRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreateCharacterReply) => void): grpc.ClientUnaryCall;
    public deleteCharacter(request: fminternal_internal_service_pb.DeleteCharacterRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DeleteCharacterReply) => void): grpc.ClientUnaryCall;
    public deleteCharacter(request: fminternal_internal_service_pb.DeleteCharacterRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DeleteCharacterReply) => void): grpc.ClientUnaryCall;
    public deleteCharacter(request: fminternal_internal_service_pb.DeleteCharacterRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DeleteCharacterReply) => void): grpc.ClientUnaryCall;
    public refreshSession(request: fminternal_internal_service_pb.RefreshSessionRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.RefreshSessionReply) => void): grpc.ClientUnaryCall;
    public refreshSession(request: fminternal_internal_service_pb.RefreshSessionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.RefreshSessionReply) => void): grpc.ClientUnaryCall;
    public refreshSession(request: fminternal_internal_service_pb.RefreshSessionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.RefreshSessionReply) => void): grpc.ClientUnaryCall;
    public logoutSession(request: fminternal_internal_service_pb.LogoutSessionRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LogoutSessionReply) => void): grpc.ClientUnaryCall;
    public logoutSession(request: fminternal_internal_service_pb.LogoutSessionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LogoutSessionReply) => void): grpc.ClientUnaryCall;
    public logoutSession(request: fminternal_internal_service_pb.LogoutSessionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LogoutSessionReply) => void): grpc.ClientUnaryCall;
    public createParty(request: fminternal_internal_service_pb.CreatePartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreatePartyReply) => void): grpc.ClientUnaryCall;
    public createParty(request: fminternal_internal_service_pb.CreatePartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreatePartyReply) => void): grpc.ClientUnaryCall;
    public createParty(request: fminternal_internal_service_pb.CreatePartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.CreatePartyReply) => void): grpc.ClientUnaryCall;
    public joinParty(request: fminternal_internal_service_pb.JoinPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.JoinPartyReply) => void): grpc.ClientUnaryCall;
    public joinParty(request: fminternal_internal_service_pb.JoinPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.JoinPartyReply) => void): grpc.ClientUnaryCall;
    public joinParty(request: fminternal_internal_service_pb.JoinPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.JoinPartyReply) => void): grpc.ClientUnaryCall;
    public leaveParty(request: fminternal_internal_service_pb.LeavePartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LeavePartyReply) => void): grpc.ClientUnaryCall;
    public leaveParty(request: fminternal_internal_service_pb.LeavePartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LeavePartyReply) => void): grpc.ClientUnaryCall;
    public leaveParty(request: fminternal_internal_service_pb.LeavePartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.LeavePartyReply) => void): grpc.ClientUnaryCall;
    public expelParty(request: fminternal_internal_service_pb.ExpelPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ExpelPartyReply) => void): grpc.ClientUnaryCall;
    public expelParty(request: fminternal_internal_service_pb.ExpelPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ExpelPartyReply) => void): grpc.ClientUnaryCall;
    public expelParty(request: fminternal_internal_service_pb.ExpelPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ExpelPartyReply) => void): grpc.ClientUnaryCall;
    public changePartyLeader(request: fminternal_internal_service_pb.ChangePartyLeaderRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ChangePartyLeaderReply) => void): grpc.ClientUnaryCall;
    public changePartyLeader(request: fminternal_internal_service_pb.ChangePartyLeaderRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ChangePartyLeaderReply) => void): grpc.ClientUnaryCall;
    public changePartyLeader(request: fminternal_internal_service_pb.ChangePartyLeaderRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.ChangePartyLeaderReply) => void): grpc.ClientUnaryCall;
    public getParty(request: fminternal_internal_service_pb.GetPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetPartyReply) => void): grpc.ClientUnaryCall;
    public getParty(request: fminternal_internal_service_pb.GetPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetPartyReply) => void): grpc.ClientUnaryCall;
    public getParty(request: fminternal_internal_service_pb.GetPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.GetPartyReply) => void): grpc.ClientUnaryCall;
    public updatePartyMember(request: fminternal_internal_service_pb.UpdatePartyMemberRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.UpdatePartyMemberReply) => void): grpc.ClientUnaryCall;
    public updatePartyMember(request: fminternal_internal_service_pb.UpdatePartyMemberRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.UpdatePartyMemberReply) => void): grpc.ClientUnaryCall;
    public updatePartyMember(request: fminternal_internal_service_pb.UpdatePartyMemberRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.UpdatePartyMemberReply) => void): grpc.ClientUnaryCall;
    public inviteParty(request: fminternal_internal_service_pb.InvitePartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.InvitePartyReply) => void): grpc.ClientUnaryCall;
    public inviteParty(request: fminternal_internal_service_pb.InvitePartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.InvitePartyReply) => void): grpc.ClientUnaryCall;
    public inviteParty(request: fminternal_internal_service_pb.InvitePartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.InvitePartyReply) => void): grpc.ClientUnaryCall;
    public denyParty(request: fminternal_internal_service_pb.DenyPartyRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DenyPartyReply) => void): grpc.ClientUnaryCall;
    public denyParty(request: fminternal_internal_service_pb.DenyPartyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DenyPartyReply) => void): grpc.ClientUnaryCall;
    public denyParty(request: fminternal_internal_service_pb.DenyPartyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.DenyPartyReply) => void): grpc.ClientUnaryCall;
    public broadcastMultiChat(request: fminternal_internal_service_pb.BroadcastMultiChatRequest, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BroadcastMultiChatReply) => void): grpc.ClientUnaryCall;
    public broadcastMultiChat(request: fminternal_internal_service_pb.BroadcastMultiChatRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BroadcastMultiChatReply) => void): grpc.ClientUnaryCall;
    public broadcastMultiChat(request: fminternal_internal_service_pb.BroadcastMultiChatRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: fminternal_internal_service_pb.BroadcastMultiChatReply) => void): grpc.ClientUnaryCall;
}
