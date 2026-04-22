// source: fminternal/internal_service.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global = (function() {
  if (this) { return this; }
  if (typeof window !== 'undefined') { return window; }
  if (typeof global !== 'undefined') { return global; }
  if (typeof self !== 'undefined') { return self; }
  return Function('return this')();
}.call(null));

goog.exportSymbol('proto.fm.internal.BeginGameTransitionReply', null, global);
goog.exportSymbol('proto.fm.internal.BeginGameTransitionRequest', null, global);
goog.exportSymbol('proto.fm.internal.ChangePartyLeaderReply', null, global);
goog.exportSymbol('proto.fm.internal.ChangePartyLeaderRequest', null, global);
goog.exportSymbol('proto.fm.internal.ChannelCatalog', null, global);
goog.exportSymbol('proto.fm.internal.CharacterOverview', null, global);
goog.exportSymbol('proto.fm.internal.CharacterPersisted', null, global);
goog.exportSymbol('proto.fm.internal.CharacterSaveEntry', null, global);
goog.exportSymbol('proto.fm.internal.CheckCharacterNameReply', null, global);
goog.exportSymbol('proto.fm.internal.CheckCharacterNameRequest', null, global);
goog.exportSymbol('proto.fm.internal.CreateCharacterReply', null, global);
goog.exportSymbol('proto.fm.internal.CreateCharacterRequest', null, global);
goog.exportSymbol('proto.fm.internal.CreatePartyReply', null, global);
goog.exportSymbol('proto.fm.internal.CreatePartyRequest', null, global);
goog.exportSymbol('proto.fm.internal.DeleteCharacterReply', null, global);
goog.exportSymbol('proto.fm.internal.DeleteCharacterRequest', null, global);
goog.exportSymbol('proto.fm.internal.DenyPartyReply', null, global);
goog.exportSymbol('proto.fm.internal.DenyPartyRequest', null, global);
goog.exportSymbol('proto.fm.internal.EnterGameReply', null, global);
goog.exportSymbol('proto.fm.internal.EnterGameRequest', null, global);
goog.exportSymbol('proto.fm.internal.ExpelPartyReply', null, global);
goog.exportSymbol('proto.fm.internal.ExpelPartyRequest', null, global);
goog.exportSymbol('proto.fm.internal.GetCharacterListReply', null, global);
goog.exportSymbol('proto.fm.internal.GetCharacterListRequest', null, global);
goog.exportSymbol('proto.fm.internal.GetPartyReply', null, global);
goog.exportSymbol('proto.fm.internal.GetPartyRequest', null, global);
goog.exportSymbol('proto.fm.internal.GetServerCatalogReply', null, global);
goog.exportSymbol('proto.fm.internal.GetServerCatalogRequest', null, global);
goog.exportSymbol('proto.fm.internal.InventoryPersisted', null, global);
goog.exportSymbol('proto.fm.internal.InvitePartyReply', null, global);
goog.exportSymbol('proto.fm.internal.InvitePartyRequest', null, global);
goog.exportSymbol('proto.fm.internal.JoinPartyReply', null, global);
goog.exportSymbol('proto.fm.internal.JoinPartyRequest', null, global);
goog.exportSymbol('proto.fm.internal.KeyLayoutBinding', null, global);
goog.exportSymbol('proto.fm.internal.LeavePartyReply', null, global);
goog.exportSymbol('proto.fm.internal.LeavePartyRequest', null, global);
goog.exportSymbol('proto.fm.internal.LoginAccountReply', null, global);
goog.exportSymbol('proto.fm.internal.LoginAccountReply.Status', null, global);
goog.exportSymbol('proto.fm.internal.LoginAccountRequest', null, global);
goog.exportSymbol('proto.fm.internal.LogoutSessionReply', null, global);
goog.exportSymbol('proto.fm.internal.LogoutSessionRequest', null, global);
goog.exportSymbol('proto.fm.internal.Party', null, global);
goog.exportSymbol('proto.fm.internal.PartyDoor', null, global);
goog.exportSymbol('proto.fm.internal.PartyErrorCode', null, global);
goog.exportSymbol('proto.fm.internal.PartyMember', null, global);
goog.exportSymbol('proto.fm.internal.PingReply', null, global);
goog.exportSymbol('proto.fm.internal.PingRequest', null, global);
goog.exportSymbol('proto.fm.internal.RefreshSessionReply', null, global);
goog.exportSymbol('proto.fm.internal.RefreshSessionRequest', null, global);
goog.exportSymbol('proto.fm.internal.SaveCharacterReply', null, global);
goog.exportSymbol('proto.fm.internal.SaveCharacterRequest', null, global);
goog.exportSymbol('proto.fm.internal.SaveCharactersReply', null, global);
goog.exportSymbol('proto.fm.internal.SaveCharactersRequest', null, global);
goog.exportSymbol('proto.fm.internal.SessionDisconnectSource', null, global);
goog.exportSymbol('proto.fm.internal.SessionErrorCode', null, global);
goog.exportSymbol('proto.fm.internal.SkillPersisted', null, global);
goog.exportSymbol('proto.fm.internal.UpdatePartyMemberReply', null, global);
goog.exportSymbol('proto.fm.internal.UpdatePartyMemberRequest', null, global);
goog.exportSymbol('proto.fm.internal.WorldCatalog', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.PingRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.PingRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.PingRequest.displayName = 'proto.fm.internal.PingRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.PingReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.PingReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.PingReply.displayName = 'proto.fm.internal.PingReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.GetServerCatalogRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.GetServerCatalogRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.GetServerCatalogRequest.displayName = 'proto.fm.internal.GetServerCatalogRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.ChannelCatalog = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.ChannelCatalog, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.ChannelCatalog.displayName = 'proto.fm.internal.ChannelCatalog';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.WorldCatalog = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.fm.internal.WorldCatalog.repeatedFields_, null);
};
goog.inherits(proto.fm.internal.WorldCatalog, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.WorldCatalog.displayName = 'proto.fm.internal.WorldCatalog';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.GetServerCatalogReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.fm.internal.GetServerCatalogReply.repeatedFields_, null);
};
goog.inherits(proto.fm.internal.GetServerCatalogReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.GetServerCatalogReply.displayName = 'proto.fm.internal.GetServerCatalogReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.EnterGameRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.EnterGameRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.EnterGameRequest.displayName = 'proto.fm.internal.EnterGameRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CharacterPersisted = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CharacterPersisted, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CharacterPersisted.displayName = 'proto.fm.internal.CharacterPersisted';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.KeyLayoutBinding = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.KeyLayoutBinding, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.KeyLayoutBinding.displayName = 'proto.fm.internal.KeyLayoutBinding';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.EnterGameReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.fm.internal.EnterGameReply.repeatedFields_, null);
};
goog.inherits(proto.fm.internal.EnterGameReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.EnterGameReply.displayName = 'proto.fm.internal.EnterGameReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.BeginGameTransitionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.BeginGameTransitionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.BeginGameTransitionRequest.displayName = 'proto.fm.internal.BeginGameTransitionRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.BeginGameTransitionReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.BeginGameTransitionReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.BeginGameTransitionReply.displayName = 'proto.fm.internal.BeginGameTransitionReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.SaveCharacterRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.SaveCharacterRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.SaveCharacterRequest.displayName = 'proto.fm.internal.SaveCharacterRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.SaveCharacterReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.SaveCharacterReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.SaveCharacterReply.displayName = 'proto.fm.internal.SaveCharacterReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CharacterSaveEntry = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.fm.internal.CharacterSaveEntry.repeatedFields_, null);
};
goog.inherits(proto.fm.internal.CharacterSaveEntry, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CharacterSaveEntry.displayName = 'proto.fm.internal.CharacterSaveEntry';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.SaveCharactersRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.fm.internal.SaveCharactersRequest.repeatedFields_, null);
};
goog.inherits(proto.fm.internal.SaveCharactersRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.SaveCharactersRequest.displayName = 'proto.fm.internal.SaveCharactersRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.SaveCharactersReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.SaveCharactersReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.SaveCharactersReply.displayName = 'proto.fm.internal.SaveCharactersReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.InventoryPersisted = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.InventoryPersisted, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.InventoryPersisted.displayName = 'proto.fm.internal.InventoryPersisted';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.SkillPersisted = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.SkillPersisted, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.SkillPersisted.displayName = 'proto.fm.internal.SkillPersisted';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.LoginAccountRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.LoginAccountRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.LoginAccountRequest.displayName = 'proto.fm.internal.LoginAccountRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.LoginAccountReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.LoginAccountReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.LoginAccountReply.displayName = 'proto.fm.internal.LoginAccountReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CharacterOverview = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CharacterOverview, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CharacterOverview.displayName = 'proto.fm.internal.CharacterOverview';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.GetCharacterListRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.GetCharacterListRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.GetCharacterListRequest.displayName = 'proto.fm.internal.GetCharacterListRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.GetCharacterListReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.fm.internal.GetCharacterListReply.repeatedFields_, null);
};
goog.inherits(proto.fm.internal.GetCharacterListReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.GetCharacterListReply.displayName = 'proto.fm.internal.GetCharacterListReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CheckCharacterNameRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CheckCharacterNameRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CheckCharacterNameRequest.displayName = 'proto.fm.internal.CheckCharacterNameRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CheckCharacterNameReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CheckCharacterNameReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CheckCharacterNameReply.displayName = 'proto.fm.internal.CheckCharacterNameReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CreateCharacterRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CreateCharacterRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CreateCharacterRequest.displayName = 'proto.fm.internal.CreateCharacterRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CreateCharacterReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CreateCharacterReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CreateCharacterReply.displayName = 'proto.fm.internal.CreateCharacterReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.DeleteCharacterRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.DeleteCharacterRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.DeleteCharacterRequest.displayName = 'proto.fm.internal.DeleteCharacterRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.DeleteCharacterReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.DeleteCharacterReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.DeleteCharacterReply.displayName = 'proto.fm.internal.DeleteCharacterReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.RefreshSessionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.RefreshSessionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.RefreshSessionRequest.displayName = 'proto.fm.internal.RefreshSessionRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.RefreshSessionReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.RefreshSessionReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.RefreshSessionReply.displayName = 'proto.fm.internal.RefreshSessionReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.LogoutSessionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.LogoutSessionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.LogoutSessionRequest.displayName = 'proto.fm.internal.LogoutSessionRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.LogoutSessionReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.LogoutSessionReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.LogoutSessionReply.displayName = 'proto.fm.internal.LogoutSessionReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.PartyDoor = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.PartyDoor, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.PartyDoor.displayName = 'proto.fm.internal.PartyDoor';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.PartyMember = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.PartyMember, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.PartyMember.displayName = 'proto.fm.internal.PartyMember';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CreatePartyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CreatePartyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CreatePartyRequest.displayName = 'proto.fm.internal.CreatePartyRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.CreatePartyReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.CreatePartyReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.CreatePartyReply.displayName = 'proto.fm.internal.CreatePartyReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.JoinPartyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.JoinPartyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.JoinPartyRequest.displayName = 'proto.fm.internal.JoinPartyRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.JoinPartyReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.JoinPartyReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.JoinPartyReply.displayName = 'proto.fm.internal.JoinPartyReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.LeavePartyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.LeavePartyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.LeavePartyRequest.displayName = 'proto.fm.internal.LeavePartyRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.LeavePartyReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.LeavePartyReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.LeavePartyReply.displayName = 'proto.fm.internal.LeavePartyReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.ExpelPartyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.ExpelPartyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.ExpelPartyRequest.displayName = 'proto.fm.internal.ExpelPartyRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.ExpelPartyReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.ExpelPartyReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.ExpelPartyReply.displayName = 'proto.fm.internal.ExpelPartyReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.ChangePartyLeaderRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.ChangePartyLeaderRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.ChangePartyLeaderRequest.displayName = 'proto.fm.internal.ChangePartyLeaderRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.ChangePartyLeaderReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.ChangePartyLeaderReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.ChangePartyLeaderReply.displayName = 'proto.fm.internal.ChangePartyLeaderReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.Party = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.fm.internal.Party.repeatedFields_, null);
};
goog.inherits(proto.fm.internal.Party, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.Party.displayName = 'proto.fm.internal.Party';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.GetPartyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.GetPartyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.GetPartyRequest.displayName = 'proto.fm.internal.GetPartyRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.GetPartyReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.GetPartyReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.GetPartyReply.displayName = 'proto.fm.internal.GetPartyReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.UpdatePartyMemberRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.UpdatePartyMemberRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.UpdatePartyMemberRequest.displayName = 'proto.fm.internal.UpdatePartyMemberRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.UpdatePartyMemberReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.UpdatePartyMemberReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.UpdatePartyMemberReply.displayName = 'proto.fm.internal.UpdatePartyMemberReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.InvitePartyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.InvitePartyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.InvitePartyRequest.displayName = 'proto.fm.internal.InvitePartyRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.InvitePartyReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.InvitePartyReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.InvitePartyReply.displayName = 'proto.fm.internal.InvitePartyReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.DenyPartyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.DenyPartyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.DenyPartyRequest.displayName = 'proto.fm.internal.DenyPartyRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.fm.internal.DenyPartyReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.fm.internal.DenyPartyReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.fm.internal.DenyPartyReply.displayName = 'proto.fm.internal.DenyPartyReply';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.PingRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.PingRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.PingRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PingRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.PingRequest}
 */
proto.fm.internal.PingRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.PingRequest;
  return proto.fm.internal.PingRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.PingRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.PingRequest}
 */
proto.fm.internal.PingRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.PingRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.PingRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.PingRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PingRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.PingReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.PingReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.PingReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PingReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    message: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.PingReply}
 */
proto.fm.internal.PingReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.PingReply;
  return proto.fm.internal.PingReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.PingReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.PingReply}
 */
proto.fm.internal.PingReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setMessage(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.PingReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.PingReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.PingReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PingReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMessage();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string message = 1;
 * @return {string}
 */
proto.fm.internal.PingReply.prototype.getMessage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.PingReply} returns this
 */
proto.fm.internal.PingReply.prototype.setMessage = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.GetServerCatalogRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.GetServerCatalogRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.GetServerCatalogRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetServerCatalogRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.GetServerCatalogRequest}
 */
proto.fm.internal.GetServerCatalogRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.GetServerCatalogRequest;
  return proto.fm.internal.GetServerCatalogRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.GetServerCatalogRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.GetServerCatalogRequest}
 */
proto.fm.internal.GetServerCatalogRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.GetServerCatalogRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.GetServerCatalogRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.GetServerCatalogRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetServerCatalogRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.ChannelCatalog.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.ChannelCatalog.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.ChannelCatalog} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ChannelCatalog.toObject = function(includeInstance, msg) {
  var f, obj = {
    channelId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    host: jspb.Message.getFieldWithDefault(msg, 2, ""),
    port: jspb.Message.getFieldWithDefault(msg, 3, 0),
    name: jspb.Message.getFieldWithDefault(msg, 4, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.ChannelCatalog}
 */
proto.fm.internal.ChannelCatalog.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.ChannelCatalog;
  return proto.fm.internal.ChannelCatalog.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.ChannelCatalog} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.ChannelCatalog}
 */
proto.fm.internal.ChannelCatalog.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setChannelId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setHost(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPort(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.ChannelCatalog.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.ChannelCatalog.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.ChannelCatalog} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ChannelCatalog.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getChannelId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getHost();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getPort();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional uint32 channel_id = 1;
 * @return {number}
 */
proto.fm.internal.ChannelCatalog.prototype.getChannelId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChannelCatalog} returns this
 */
proto.fm.internal.ChannelCatalog.prototype.setChannelId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string host = 2;
 * @return {string}
 */
proto.fm.internal.ChannelCatalog.prototype.getHost = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.ChannelCatalog} returns this
 */
proto.fm.internal.ChannelCatalog.prototype.setHost = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional uint32 port = 3;
 * @return {number}
 */
proto.fm.internal.ChannelCatalog.prototype.getPort = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChannelCatalog} returns this
 */
proto.fm.internal.ChannelCatalog.prototype.setPort = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional string name = 4;
 * @return {string}
 */
proto.fm.internal.ChannelCatalog.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.ChannelCatalog} returns this
 */
proto.fm.internal.ChannelCatalog.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.fm.internal.WorldCatalog.repeatedFields_ = [5];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.WorldCatalog.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.WorldCatalog.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.WorldCatalog} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.WorldCatalog.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    worldName: jspb.Message.getFieldWithDefault(msg, 2, ""),
    flag: jspb.Message.getFieldWithDefault(msg, 3, 0),
    eventMessage: jspb.Message.getFieldWithDefault(msg, 4, ""),
    channelsList: jspb.Message.toObjectList(msg.getChannelsList(),
    proto.fm.internal.ChannelCatalog.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.WorldCatalog}
 */
proto.fm.internal.WorldCatalog.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.WorldCatalog;
  return proto.fm.internal.WorldCatalog.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.WorldCatalog} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.WorldCatalog}
 */
proto.fm.internal.WorldCatalog.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setWorldName(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setFlag(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setEventMessage(value);
      break;
    case 5:
      var value = new proto.fm.internal.ChannelCatalog;
      reader.readMessage(value,proto.fm.internal.ChannelCatalog.deserializeBinaryFromReader);
      msg.addChannels(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.WorldCatalog.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.WorldCatalog.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.WorldCatalog} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.WorldCatalog.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getWorldName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getFlag();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getEventMessage();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getChannelsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.fm.internal.ChannelCatalog.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.WorldCatalog.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.WorldCatalog} returns this
 */
proto.fm.internal.WorldCatalog.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string world_name = 2;
 * @return {string}
 */
proto.fm.internal.WorldCatalog.prototype.getWorldName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.WorldCatalog} returns this
 */
proto.fm.internal.WorldCatalog.prototype.setWorldName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional uint32 flag = 3;
 * @return {number}
 */
proto.fm.internal.WorldCatalog.prototype.getFlag = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.WorldCatalog} returns this
 */
proto.fm.internal.WorldCatalog.prototype.setFlag = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional string event_message = 4;
 * @return {string}
 */
proto.fm.internal.WorldCatalog.prototype.getEventMessage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.WorldCatalog} returns this
 */
proto.fm.internal.WorldCatalog.prototype.setEventMessage = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * repeated ChannelCatalog channels = 5;
 * @return {!Array<!proto.fm.internal.ChannelCatalog>}
 */
proto.fm.internal.WorldCatalog.prototype.getChannelsList = function() {
  return /** @type{!Array<!proto.fm.internal.ChannelCatalog>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.ChannelCatalog, 5));
};


/**
 * @param {!Array<!proto.fm.internal.ChannelCatalog>} value
 * @return {!proto.fm.internal.WorldCatalog} returns this
*/
proto.fm.internal.WorldCatalog.prototype.setChannelsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.fm.internal.ChannelCatalog=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.ChannelCatalog}
 */
proto.fm.internal.WorldCatalog.prototype.addChannels = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.fm.internal.ChannelCatalog, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.WorldCatalog} returns this
 */
proto.fm.internal.WorldCatalog.prototype.clearChannelsList = function() {
  return this.setChannelsList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.fm.internal.GetServerCatalogReply.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.GetServerCatalogReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.GetServerCatalogReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.GetServerCatalogReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetServerCatalogReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldsList: jspb.Message.toObjectList(msg.getWorldsList(),
    proto.fm.internal.WorldCatalog.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.GetServerCatalogReply}
 */
proto.fm.internal.GetServerCatalogReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.GetServerCatalogReply;
  return proto.fm.internal.GetServerCatalogReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.GetServerCatalogReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.GetServerCatalogReply}
 */
proto.fm.internal.GetServerCatalogReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.fm.internal.WorldCatalog;
      reader.readMessage(value,proto.fm.internal.WorldCatalog.deserializeBinaryFromReader);
      msg.addWorlds(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.GetServerCatalogReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.GetServerCatalogReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.GetServerCatalogReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetServerCatalogReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.fm.internal.WorldCatalog.serializeBinaryToWriter
    );
  }
};


/**
 * repeated WorldCatalog worlds = 1;
 * @return {!Array<!proto.fm.internal.WorldCatalog>}
 */
proto.fm.internal.GetServerCatalogReply.prototype.getWorldsList = function() {
  return /** @type{!Array<!proto.fm.internal.WorldCatalog>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.WorldCatalog, 1));
};


/**
 * @param {!Array<!proto.fm.internal.WorldCatalog>} value
 * @return {!proto.fm.internal.GetServerCatalogReply} returns this
*/
proto.fm.internal.GetServerCatalogReply.prototype.setWorldsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.fm.internal.WorldCatalog=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.WorldCatalog}
 */
proto.fm.internal.GetServerCatalogReply.prototype.addWorlds = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.fm.internal.WorldCatalog, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.GetServerCatalogReply} returns this
 */
proto.fm.internal.GetServerCatalogReply.prototype.clearWorldsList = function() {
  return this.setWorldsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.EnterGameRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.EnterGameRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.EnterGameRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.EnterGameRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    characterId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    channelId: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.EnterGameRequest}
 */
proto.fm.internal.EnterGameRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.EnterGameRequest;
  return proto.fm.internal.EnterGameRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.EnterGameRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.EnterGameRequest}
 */
proto.fm.internal.EnterGameRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setChannelId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.EnterGameRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.EnterGameRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.EnterGameRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.EnterGameRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getChannelId();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.EnterGameRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.EnterGameRequest} returns this
 */
proto.fm.internal.EnterGameRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 character_id = 2;
 * @return {number}
 */
proto.fm.internal.EnterGameRequest.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.EnterGameRequest} returns this
 */
proto.fm.internal.EnterGameRequest.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint32 channel_id = 3;
 * @return {number}
 */
proto.fm.internal.EnterGameRequest.prototype.getChannelId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.EnterGameRequest} returns this
 */
proto.fm.internal.EnterGameRequest.prototype.setChannelId = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CharacterPersisted.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CharacterPersisted.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CharacterPersisted} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CharacterPersisted.toObject = function(includeInstance, msg) {
  var f, obj = {
    characterId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    worldId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    name: jspb.Message.getFieldWithDefault(msg, 3, ""),
    gender: jspb.Message.getFieldWithDefault(msg, 4, 0),
    skinColor: jspb.Message.getFieldWithDefault(msg, 5, 0),
    face: jspb.Message.getFieldWithDefault(msg, 6, 0),
    hair: jspb.Message.getFieldWithDefault(msg, 7, 0),
    level: jspb.Message.getFieldWithDefault(msg, 8, 0),
    classId: jspb.Message.getFieldWithDefault(msg, 9, 0),
    str: jspb.Message.getFieldWithDefault(msg, 10, 0),
    dex: jspb.Message.getFieldWithDefault(msg, 11, 0),
    intStat: jspb.Message.getFieldWithDefault(msg, 12, 0),
    luk: jspb.Message.getFieldWithDefault(msg, 13, 0),
    hp: jspb.Message.getFieldWithDefault(msg, 14, 0),
    maxHp: jspb.Message.getFieldWithDefault(msg, 15, 0),
    mp: jspb.Message.getFieldWithDefault(msg, 16, 0),
    maxMp: jspb.Message.getFieldWithDefault(msg, 17, 0),
    abilityPoint: jspb.Message.getFieldWithDefault(msg, 18, 0),
    exp: jspb.Message.getFieldWithDefault(msg, 19, 0),
    mapId: jspb.Message.getFieldWithDefault(msg, 20, 0),
    spawnPoint: jspb.Message.getFieldWithDefault(msg, 21, 0),
    positionX: jspb.Message.getFieldWithDefault(msg, 22, 0),
    positionY: jspb.Message.getFieldWithDefault(msg, 23, 0),
    stance: jspb.Message.getFieldWithDefault(msg, 24, 0),
    meso: jspb.Message.getFieldWithDefault(msg, 25, 0),
    skillPoint: jspb.Message.getFieldWithDefault(msg, 26, 0),
    accountId: jspb.Message.getFieldWithDefault(msg, 27, 0),
    role: jspb.Message.getFieldWithDefault(msg, 28, 0),
    hidden: jspb.Message.getBooleanFieldWithDefault(msg, 29, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CharacterPersisted}
 */
proto.fm.internal.CharacterPersisted.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CharacterPersisted;
  return proto.fm.internal.CharacterPersisted.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CharacterPersisted} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CharacterPersisted}
 */
proto.fm.internal.CharacterPersisted.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setGender(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSkinColor(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setFace(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setHair(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setLevel(value);
      break;
    case 9:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setClassId(value);
      break;
    case 10:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setStr(value);
      break;
    case 11:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setDex(value);
      break;
    case 12:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setIntStat(value);
      break;
    case 13:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setLuk(value);
      break;
    case 14:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setHp(value);
      break;
    case 15:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMaxHp(value);
      break;
    case 16:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMp(value);
      break;
    case 17:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMaxMp(value);
      break;
    case 18:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAbilityPoint(value);
      break;
    case 19:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setExp(value);
      break;
    case 20:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMapId(value);
      break;
    case 21:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSpawnPoint(value);
      break;
    case 22:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPositionX(value);
      break;
    case 23:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPositionY(value);
      break;
    case 24:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setStance(value);
      break;
    case 25:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setMeso(value);
      break;
    case 26:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSkillPoint(value);
      break;
    case 27:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 28:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setRole(value);
      break;
    case 29:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setHidden(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CharacterPersisted.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CharacterPersisted.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CharacterPersisted} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CharacterPersisted.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getGender();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
  f = message.getSkinColor();
  if (f !== 0) {
    writer.writeUint32(
      5,
      f
    );
  }
  f = message.getFace();
  if (f !== 0) {
    writer.writeUint32(
      6,
      f
    );
  }
  f = message.getHair();
  if (f !== 0) {
    writer.writeUint32(
      7,
      f
    );
  }
  f = message.getLevel();
  if (f !== 0) {
    writer.writeUint32(
      8,
      f
    );
  }
  f = message.getClassId();
  if (f !== 0) {
    writer.writeUint32(
      9,
      f
    );
  }
  f = message.getStr();
  if (f !== 0) {
    writer.writeUint32(
      10,
      f
    );
  }
  f = message.getDex();
  if (f !== 0) {
    writer.writeUint32(
      11,
      f
    );
  }
  f = message.getIntStat();
  if (f !== 0) {
    writer.writeUint32(
      12,
      f
    );
  }
  f = message.getLuk();
  if (f !== 0) {
    writer.writeUint32(
      13,
      f
    );
  }
  f = message.getHp();
  if (f !== 0) {
    writer.writeUint32(
      14,
      f
    );
  }
  f = message.getMaxHp();
  if (f !== 0) {
    writer.writeUint32(
      15,
      f
    );
  }
  f = message.getMp();
  if (f !== 0) {
    writer.writeUint32(
      16,
      f
    );
  }
  f = message.getMaxMp();
  if (f !== 0) {
    writer.writeUint32(
      17,
      f
    );
  }
  f = message.getAbilityPoint();
  if (f !== 0) {
    writer.writeUint32(
      18,
      f
    );
  }
  f = message.getExp();
  if (f !== 0) {
    writer.writeUint32(
      19,
      f
    );
  }
  f = message.getMapId();
  if (f !== 0) {
    writer.writeUint32(
      20,
      f
    );
  }
  f = message.getSpawnPoint();
  if (f !== 0) {
    writer.writeUint32(
      21,
      f
    );
  }
  f = message.getPositionX();
  if (f !== 0) {
    writer.writeInt32(
      22,
      f
    );
  }
  f = message.getPositionY();
  if (f !== 0) {
    writer.writeInt32(
      23,
      f
    );
  }
  f = message.getStance();
  if (f !== 0) {
    writer.writeUint32(
      24,
      f
    );
  }
  f = message.getMeso();
  if (f !== 0) {
    writer.writeInt32(
      25,
      f
    );
  }
  f = message.getSkillPoint();
  if (f !== 0) {
    writer.writeUint32(
      26,
      f
    );
  }
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      27,
      f
    );
  }
  f = message.getRole();
  if (f !== 0) {
    writer.writeUint32(
      28,
      f
    );
  }
  f = message.getHidden();
  if (f) {
    writer.writeBool(
      29,
      f
    );
  }
};


/**
 * optional uint32 character_id = 1;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 world_id = 2;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.fm.internal.CharacterPersisted.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional uint32 gender = 4;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getGender = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setGender = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional uint32 skin_color = 5;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getSkinColor = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setSkinColor = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional uint32 face = 6;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getFace = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setFace = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional uint32 hair = 7;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getHair = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setHair = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional uint32 level = 8;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getLevel = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setLevel = function(value) {
  return jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * optional uint32 class_id = 9;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getClassId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setClassId = function(value) {
  return jspb.Message.setProto3IntField(this, 9, value);
};


/**
 * optional uint32 str = 10;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getStr = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 10, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setStr = function(value) {
  return jspb.Message.setProto3IntField(this, 10, value);
};


/**
 * optional uint32 dex = 11;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getDex = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setDex = function(value) {
  return jspb.Message.setProto3IntField(this, 11, value);
};


/**
 * optional uint32 int_stat = 12;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getIntStat = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 12, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setIntStat = function(value) {
  return jspb.Message.setProto3IntField(this, 12, value);
};


/**
 * optional uint32 luk = 13;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getLuk = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 13, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setLuk = function(value) {
  return jspb.Message.setProto3IntField(this, 13, value);
};


/**
 * optional uint32 hp = 14;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getHp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 14, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setHp = function(value) {
  return jspb.Message.setProto3IntField(this, 14, value);
};


/**
 * optional uint32 max_hp = 15;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getMaxHp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 15, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setMaxHp = function(value) {
  return jspb.Message.setProto3IntField(this, 15, value);
};


/**
 * optional uint32 mp = 16;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getMp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 16, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setMp = function(value) {
  return jspb.Message.setProto3IntField(this, 16, value);
};


/**
 * optional uint32 max_mp = 17;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getMaxMp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 17, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setMaxMp = function(value) {
  return jspb.Message.setProto3IntField(this, 17, value);
};


/**
 * optional uint32 ability_point = 18;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getAbilityPoint = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 18, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setAbilityPoint = function(value) {
  return jspb.Message.setProto3IntField(this, 18, value);
};


/**
 * optional uint32 exp = 19;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getExp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 19, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setExp = function(value) {
  return jspb.Message.setProto3IntField(this, 19, value);
};


/**
 * optional uint32 map_id = 20;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getMapId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 20, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setMapId = function(value) {
  return jspb.Message.setProto3IntField(this, 20, value);
};


/**
 * optional uint32 spawn_point = 21;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getSpawnPoint = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 21, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setSpawnPoint = function(value) {
  return jspb.Message.setProto3IntField(this, 21, value);
};


/**
 * optional int32 position_x = 22;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getPositionX = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 22, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setPositionX = function(value) {
  return jspb.Message.setProto3IntField(this, 22, value);
};


/**
 * optional int32 position_y = 23;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getPositionY = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 23, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setPositionY = function(value) {
  return jspb.Message.setProto3IntField(this, 23, value);
};


/**
 * optional uint32 stance = 24;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getStance = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 24, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setStance = function(value) {
  return jspb.Message.setProto3IntField(this, 24, value);
};


/**
 * optional int32 meso = 25;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getMeso = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 25, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setMeso = function(value) {
  return jspb.Message.setProto3IntField(this, 25, value);
};


/**
 * optional uint32 skill_point = 26;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getSkillPoint = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 26, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setSkillPoint = function(value) {
  return jspb.Message.setProto3IntField(this, 26, value);
};


/**
 * optional uint32 account_id = 27;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 27, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 27, value);
};


/**
 * optional uint32 role = 28;
 * @return {number}
 */
proto.fm.internal.CharacterPersisted.prototype.getRole = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 28, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setRole = function(value) {
  return jspb.Message.setProto3IntField(this, 28, value);
};


/**
 * optional bool hidden = 29;
 * @return {boolean}
 */
proto.fm.internal.CharacterPersisted.prototype.getHidden = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 29, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.CharacterPersisted} returns this
 */
proto.fm.internal.CharacterPersisted.prototype.setHidden = function(value) {
  return jspb.Message.setProto3BooleanField(this, 29, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.KeyLayoutBinding.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.KeyLayoutBinding.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.KeyLayoutBinding} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.KeyLayoutBinding.toObject = function(includeInstance, msg) {
  var f, obj = {
    slot: jspb.Message.getFieldWithDefault(msg, 1, 0),
    type: jspb.Message.getFieldWithDefault(msg, 2, 0),
    action: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.KeyLayoutBinding}
 */
proto.fm.internal.KeyLayoutBinding.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.KeyLayoutBinding;
  return proto.fm.internal.KeyLayoutBinding.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.KeyLayoutBinding} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.KeyLayoutBinding}
 */
proto.fm.internal.KeyLayoutBinding.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setSlot(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setType(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setAction(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.KeyLayoutBinding.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.KeyLayoutBinding.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.KeyLayoutBinding} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.KeyLayoutBinding.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSlot();
  if (f !== 0) {
    writer.writeInt32(
      1,
      f
    );
  }
  f = message.getType();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getAction();
  if (f !== 0) {
    writer.writeInt32(
      3,
      f
    );
  }
};


/**
 * optional int32 slot = 1;
 * @return {number}
 */
proto.fm.internal.KeyLayoutBinding.prototype.getSlot = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.KeyLayoutBinding} returns this
 */
proto.fm.internal.KeyLayoutBinding.prototype.setSlot = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 type = 2;
 * @return {number}
 */
proto.fm.internal.KeyLayoutBinding.prototype.getType = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.KeyLayoutBinding} returns this
 */
proto.fm.internal.KeyLayoutBinding.prototype.setType = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int32 action = 3;
 * @return {number}
 */
proto.fm.internal.KeyLayoutBinding.prototype.getAction = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.KeyLayoutBinding} returns this
 */
proto.fm.internal.KeyLayoutBinding.prototype.setAction = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.fm.internal.EnterGameReply.repeatedFields_ = [3,4,5];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.EnterGameReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.EnterGameReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.EnterGameReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.EnterGameReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    found: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    character: (f = msg.getCharacter()) && proto.fm.internal.CharacterPersisted.toObject(includeInstance, f),
    inventoryList: jspb.Message.toObjectList(msg.getInventoryList(),
    proto.fm.internal.InventoryPersisted.toObject, includeInstance),
    skillsList: jspb.Message.toObjectList(msg.getSkillsList(),
    proto.fm.internal.SkillPersisted.toObject, includeInstance),
    keyLayoutList: jspb.Message.toObjectList(msg.getKeyLayoutList(),
    proto.fm.internal.KeyLayoutBinding.toObject, includeInstance),
    partyId: jspb.Message.getFieldWithDefault(msg, 6, 0),
    guildId: jspb.Message.getFieldWithDefault(msg, 7, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.EnterGameReply}
 */
proto.fm.internal.EnterGameReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.EnterGameReply;
  return proto.fm.internal.EnterGameReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.EnterGameReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.EnterGameReply}
 */
proto.fm.internal.EnterGameReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setFound(value);
      break;
    case 2:
      var value = new proto.fm.internal.CharacterPersisted;
      reader.readMessage(value,proto.fm.internal.CharacterPersisted.deserializeBinaryFromReader);
      msg.setCharacter(value);
      break;
    case 3:
      var value = new proto.fm.internal.InventoryPersisted;
      reader.readMessage(value,proto.fm.internal.InventoryPersisted.deserializeBinaryFromReader);
      msg.addInventory(value);
      break;
    case 4:
      var value = new proto.fm.internal.SkillPersisted;
      reader.readMessage(value,proto.fm.internal.SkillPersisted.deserializeBinaryFromReader);
      msg.addSkills(value);
      break;
    case 5:
      var value = new proto.fm.internal.KeyLayoutBinding;
      reader.readMessage(value,proto.fm.internal.KeyLayoutBinding.deserializeBinaryFromReader);
      msg.addKeyLayout(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setGuildId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.EnterGameReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.EnterGameReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.EnterGameReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.EnterGameReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFound();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getCharacter();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.fm.internal.CharacterPersisted.serializeBinaryToWriter
    );
  }
  f = message.getInventoryList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.fm.internal.InventoryPersisted.serializeBinaryToWriter
    );
  }
  f = message.getSkillsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.fm.internal.SkillPersisted.serializeBinaryToWriter
    );
  }
  f = message.getKeyLayoutList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.fm.internal.KeyLayoutBinding.serializeBinaryToWriter
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 6));
  if (f != null) {
    writer.writeUint32(
      6,
      f
    );
  }
  f = message.getGuildId();
  if (f !== 0) {
    writer.writeUint32(
      7,
      f
    );
  }
};


/**
 * optional bool found = 1;
 * @return {boolean}
 */
proto.fm.internal.EnterGameReply.prototype.getFound = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.setFound = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional CharacterPersisted character = 2;
 * @return {?proto.fm.internal.CharacterPersisted}
 */
proto.fm.internal.EnterGameReply.prototype.getCharacter = function() {
  return /** @type{?proto.fm.internal.CharacterPersisted} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.CharacterPersisted, 2));
};


/**
 * @param {?proto.fm.internal.CharacterPersisted|undefined} value
 * @return {!proto.fm.internal.EnterGameReply} returns this
*/
proto.fm.internal.EnterGameReply.prototype.setCharacter = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.clearCharacter = function() {
  return this.setCharacter(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.EnterGameReply.prototype.hasCharacter = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * repeated InventoryPersisted inventory = 3;
 * @return {!Array<!proto.fm.internal.InventoryPersisted>}
 */
proto.fm.internal.EnterGameReply.prototype.getInventoryList = function() {
  return /** @type{!Array<!proto.fm.internal.InventoryPersisted>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.InventoryPersisted, 3));
};


/**
 * @param {!Array<!proto.fm.internal.InventoryPersisted>} value
 * @return {!proto.fm.internal.EnterGameReply} returns this
*/
proto.fm.internal.EnterGameReply.prototype.setInventoryList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.fm.internal.InventoryPersisted=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.InventoryPersisted}
 */
proto.fm.internal.EnterGameReply.prototype.addInventory = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.fm.internal.InventoryPersisted, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.clearInventoryList = function() {
  return this.setInventoryList([]);
};


/**
 * repeated SkillPersisted skills = 4;
 * @return {!Array<!proto.fm.internal.SkillPersisted>}
 */
proto.fm.internal.EnterGameReply.prototype.getSkillsList = function() {
  return /** @type{!Array<!proto.fm.internal.SkillPersisted>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.SkillPersisted, 4));
};


/**
 * @param {!Array<!proto.fm.internal.SkillPersisted>} value
 * @return {!proto.fm.internal.EnterGameReply} returns this
*/
proto.fm.internal.EnterGameReply.prototype.setSkillsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.fm.internal.SkillPersisted=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.SkillPersisted}
 */
proto.fm.internal.EnterGameReply.prototype.addSkills = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.fm.internal.SkillPersisted, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.clearSkillsList = function() {
  return this.setSkillsList([]);
};


/**
 * repeated KeyLayoutBinding key_layout = 5;
 * @return {!Array<!proto.fm.internal.KeyLayoutBinding>}
 */
proto.fm.internal.EnterGameReply.prototype.getKeyLayoutList = function() {
  return /** @type{!Array<!proto.fm.internal.KeyLayoutBinding>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.KeyLayoutBinding, 5));
};


/**
 * @param {!Array<!proto.fm.internal.KeyLayoutBinding>} value
 * @return {!proto.fm.internal.EnterGameReply} returns this
*/
proto.fm.internal.EnterGameReply.prototype.setKeyLayoutList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.fm.internal.KeyLayoutBinding=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.KeyLayoutBinding}
 */
proto.fm.internal.EnterGameReply.prototype.addKeyLayout = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.fm.internal.KeyLayoutBinding, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.clearKeyLayoutList = function() {
  return this.setKeyLayoutList([]);
};


/**
 * optional uint32 party_id = 6;
 * @return {number}
 */
proto.fm.internal.EnterGameReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 6, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 6, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.EnterGameReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional uint32 guild_id = 7;
 * @return {number}
 */
proto.fm.internal.EnterGameReply.prototype.getGuildId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.EnterGameReply} returns this
 */
proto.fm.internal.EnterGameReply.prototype.setGuildId = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.BeginGameTransitionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.BeginGameTransitionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.BeginGameTransitionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    accountId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    characterId: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.BeginGameTransitionRequest}
 */
proto.fm.internal.BeginGameTransitionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.BeginGameTransitionRequest;
  return proto.fm.internal.BeginGameTransitionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.BeginGameTransitionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.BeginGameTransitionRequest}
 */
proto.fm.internal.BeginGameTransitionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.BeginGameTransitionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.BeginGameTransitionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.BeginGameTransitionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.BeginGameTransitionRequest} returns this
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 account_id = 2;
 * @return {number}
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.BeginGameTransitionRequest} returns this
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint32 character_id = 3;
 * @return {number}
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.BeginGameTransitionRequest} returns this
 */
proto.fm.internal.BeginGameTransitionRequest.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.BeginGameTransitionReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.BeginGameTransitionReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.BeginGameTransitionReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.BeginGameTransitionReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.BeginGameTransitionReply}
 */
proto.fm.internal.BeginGameTransitionReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.BeginGameTransitionReply;
  return proto.fm.internal.BeginGameTransitionReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.BeginGameTransitionReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.BeginGameTransitionReply}
 */
proto.fm.internal.BeginGameTransitionReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.SessionErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.BeginGameTransitionReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.BeginGameTransitionReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.BeginGameTransitionReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.BeginGameTransitionReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.BeginGameTransitionReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.BeginGameTransitionReply} returns this
 */
proto.fm.internal.BeginGameTransitionReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional SessionErrorCode error_code = 2;
 * @return {!proto.fm.internal.SessionErrorCode}
 */
proto.fm.internal.BeginGameTransitionReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.SessionErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.SessionErrorCode} value
 * @return {!proto.fm.internal.BeginGameTransitionReply} returns this
 */
proto.fm.internal.BeginGameTransitionReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.SaveCharacterRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.SaveCharacterRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.SaveCharacterRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharacterRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    character: (f = msg.getCharacter()) && proto.fm.internal.CharacterPersisted.toObject(includeInstance, f),
    baseLooksMap: (f = msg.getBaseLooksMap()) ? f.toObject(includeInstance, undefined) : [],
    overlaysMap: (f = msg.getOverlaysMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.SaveCharacterRequest}
 */
proto.fm.internal.SaveCharacterRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.SaveCharacterRequest;
  return proto.fm.internal.SaveCharacterRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.SaveCharacterRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.SaveCharacterRequest}
 */
proto.fm.internal.SaveCharacterRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.fm.internal.CharacterPersisted;
      reader.readMessage(value,proto.fm.internal.CharacterPersisted.deserializeBinaryFromReader);
      msg.setCharacter(value);
      break;
    case 2:
      var value = msg.getBaseLooksMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readInt32, jspb.BinaryReader.prototype.readUint32, null, 0, 0);
         });
      break;
    case 3:
      var value = msg.getOverlaysMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readInt32, jspb.BinaryReader.prototype.readUint32, null, 0, 0);
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.SaveCharacterRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.SaveCharacterRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.SaveCharacterRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharacterRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCharacter();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.fm.internal.CharacterPersisted.serializeBinaryToWriter
    );
  }
  f = message.getBaseLooksMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(2, writer, jspb.BinaryWriter.prototype.writeInt32, jspb.BinaryWriter.prototype.writeUint32);
  }
  f = message.getOverlaysMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(3, writer, jspb.BinaryWriter.prototype.writeInt32, jspb.BinaryWriter.prototype.writeUint32);
  }
};


/**
 * optional CharacterPersisted character = 1;
 * @return {?proto.fm.internal.CharacterPersisted}
 */
proto.fm.internal.SaveCharacterRequest.prototype.getCharacter = function() {
  return /** @type{?proto.fm.internal.CharacterPersisted} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.CharacterPersisted, 1));
};


/**
 * @param {?proto.fm.internal.CharacterPersisted|undefined} value
 * @return {!proto.fm.internal.SaveCharacterRequest} returns this
*/
proto.fm.internal.SaveCharacterRequest.prototype.setCharacter = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.SaveCharacterRequest} returns this
 */
proto.fm.internal.SaveCharacterRequest.prototype.clearCharacter = function() {
  return this.setCharacter(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.SaveCharacterRequest.prototype.hasCharacter = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * map<int32, uint32> base_looks = 2;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<number,number>}
 */
proto.fm.internal.SaveCharacterRequest.prototype.getBaseLooksMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<number,number>} */ (
      jspb.Message.getMapField(this, 2, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.fm.internal.SaveCharacterRequest} returns this
 */
proto.fm.internal.SaveCharacterRequest.prototype.clearBaseLooksMap = function() {
  this.getBaseLooksMap().clear();
  return this;};


/**
 * map<int32, uint32> overlays = 3;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<number,number>}
 */
proto.fm.internal.SaveCharacterRequest.prototype.getOverlaysMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<number,number>} */ (
      jspb.Message.getMapField(this, 3, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.fm.internal.SaveCharacterRequest} returns this
 */
proto.fm.internal.SaveCharacterRequest.prototype.clearOverlaysMap = function() {
  this.getOverlaysMap().clear();
  return this;};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.SaveCharacterReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.SaveCharacterReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.SaveCharacterReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharacterReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.SaveCharacterReply}
 */
proto.fm.internal.SaveCharacterReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.SaveCharacterReply;
  return proto.fm.internal.SaveCharacterReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.SaveCharacterReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.SaveCharacterReply}
 */
proto.fm.internal.SaveCharacterReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.SaveCharacterReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.SaveCharacterReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.SaveCharacterReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharacterReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.SaveCharacterReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.SaveCharacterReply} returns this
 */
proto.fm.internal.SaveCharacterReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.fm.internal.CharacterSaveEntry.repeatedFields_ = [4,5,6];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CharacterSaveEntry.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CharacterSaveEntry.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CharacterSaveEntry} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CharacterSaveEntry.toObject = function(includeInstance, msg) {
  var f, obj = {
    character: (f = msg.getCharacter()) && proto.fm.internal.CharacterPersisted.toObject(includeInstance, f),
    baseLooksMap: (f = msg.getBaseLooksMap()) ? f.toObject(includeInstance, undefined) : [],
    overlaysMap: (f = msg.getOverlaysMap()) ? f.toObject(includeInstance, undefined) : [],
    inventoryList: jspb.Message.toObjectList(msg.getInventoryList(),
    proto.fm.internal.InventoryPersisted.toObject, includeInstance),
    skillsList: jspb.Message.toObjectList(msg.getSkillsList(),
    proto.fm.internal.SkillPersisted.toObject, includeInstance),
    keyLayoutList: jspb.Message.toObjectList(msg.getKeyLayoutList(),
    proto.fm.internal.KeyLayoutBinding.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CharacterSaveEntry}
 */
proto.fm.internal.CharacterSaveEntry.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CharacterSaveEntry;
  return proto.fm.internal.CharacterSaveEntry.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CharacterSaveEntry} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CharacterSaveEntry}
 */
proto.fm.internal.CharacterSaveEntry.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.fm.internal.CharacterPersisted;
      reader.readMessage(value,proto.fm.internal.CharacterPersisted.deserializeBinaryFromReader);
      msg.setCharacter(value);
      break;
    case 2:
      var value = msg.getBaseLooksMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readInt32, jspb.BinaryReader.prototype.readUint32, null, 0, 0);
         });
      break;
    case 3:
      var value = msg.getOverlaysMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readInt32, jspb.BinaryReader.prototype.readUint32, null, 0, 0);
         });
      break;
    case 4:
      var value = new proto.fm.internal.InventoryPersisted;
      reader.readMessage(value,proto.fm.internal.InventoryPersisted.deserializeBinaryFromReader);
      msg.addInventory(value);
      break;
    case 5:
      var value = new proto.fm.internal.SkillPersisted;
      reader.readMessage(value,proto.fm.internal.SkillPersisted.deserializeBinaryFromReader);
      msg.addSkills(value);
      break;
    case 6:
      var value = new proto.fm.internal.KeyLayoutBinding;
      reader.readMessage(value,proto.fm.internal.KeyLayoutBinding.deserializeBinaryFromReader);
      msg.addKeyLayout(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CharacterSaveEntry.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CharacterSaveEntry.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CharacterSaveEntry} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CharacterSaveEntry.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCharacter();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.fm.internal.CharacterPersisted.serializeBinaryToWriter
    );
  }
  f = message.getBaseLooksMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(2, writer, jspb.BinaryWriter.prototype.writeInt32, jspb.BinaryWriter.prototype.writeUint32);
  }
  f = message.getOverlaysMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(3, writer, jspb.BinaryWriter.prototype.writeInt32, jspb.BinaryWriter.prototype.writeUint32);
  }
  f = message.getInventoryList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.fm.internal.InventoryPersisted.serializeBinaryToWriter
    );
  }
  f = message.getSkillsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.fm.internal.SkillPersisted.serializeBinaryToWriter
    );
  }
  f = message.getKeyLayoutList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      6,
      f,
      proto.fm.internal.KeyLayoutBinding.serializeBinaryToWriter
    );
  }
};


/**
 * optional CharacterPersisted character = 1;
 * @return {?proto.fm.internal.CharacterPersisted}
 */
proto.fm.internal.CharacterSaveEntry.prototype.getCharacter = function() {
  return /** @type{?proto.fm.internal.CharacterPersisted} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.CharacterPersisted, 1));
};


/**
 * @param {?proto.fm.internal.CharacterPersisted|undefined} value
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
*/
proto.fm.internal.CharacterSaveEntry.prototype.setCharacter = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
 */
proto.fm.internal.CharacterSaveEntry.prototype.clearCharacter = function() {
  return this.setCharacter(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.CharacterSaveEntry.prototype.hasCharacter = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * map<int32, uint32> base_looks = 2;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<number,number>}
 */
proto.fm.internal.CharacterSaveEntry.prototype.getBaseLooksMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<number,number>} */ (
      jspb.Message.getMapField(this, 2, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
 */
proto.fm.internal.CharacterSaveEntry.prototype.clearBaseLooksMap = function() {
  this.getBaseLooksMap().clear();
  return this;};


/**
 * map<int32, uint32> overlays = 3;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<number,number>}
 */
proto.fm.internal.CharacterSaveEntry.prototype.getOverlaysMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<number,number>} */ (
      jspb.Message.getMapField(this, 3, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
 */
proto.fm.internal.CharacterSaveEntry.prototype.clearOverlaysMap = function() {
  this.getOverlaysMap().clear();
  return this;};


/**
 * repeated InventoryPersisted inventory = 4;
 * @return {!Array<!proto.fm.internal.InventoryPersisted>}
 */
proto.fm.internal.CharacterSaveEntry.prototype.getInventoryList = function() {
  return /** @type{!Array<!proto.fm.internal.InventoryPersisted>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.InventoryPersisted, 4));
};


/**
 * @param {!Array<!proto.fm.internal.InventoryPersisted>} value
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
*/
proto.fm.internal.CharacterSaveEntry.prototype.setInventoryList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.fm.internal.InventoryPersisted=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.InventoryPersisted}
 */
proto.fm.internal.CharacterSaveEntry.prototype.addInventory = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.fm.internal.InventoryPersisted, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
 */
proto.fm.internal.CharacterSaveEntry.prototype.clearInventoryList = function() {
  return this.setInventoryList([]);
};


/**
 * repeated SkillPersisted skills = 5;
 * @return {!Array<!proto.fm.internal.SkillPersisted>}
 */
proto.fm.internal.CharacterSaveEntry.prototype.getSkillsList = function() {
  return /** @type{!Array<!proto.fm.internal.SkillPersisted>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.SkillPersisted, 5));
};


/**
 * @param {!Array<!proto.fm.internal.SkillPersisted>} value
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
*/
proto.fm.internal.CharacterSaveEntry.prototype.setSkillsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.fm.internal.SkillPersisted=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.SkillPersisted}
 */
proto.fm.internal.CharacterSaveEntry.prototype.addSkills = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.fm.internal.SkillPersisted, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
 */
proto.fm.internal.CharacterSaveEntry.prototype.clearSkillsList = function() {
  return this.setSkillsList([]);
};


/**
 * repeated KeyLayoutBinding key_layout = 6;
 * @return {!Array<!proto.fm.internal.KeyLayoutBinding>}
 */
proto.fm.internal.CharacterSaveEntry.prototype.getKeyLayoutList = function() {
  return /** @type{!Array<!proto.fm.internal.KeyLayoutBinding>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.KeyLayoutBinding, 6));
};


/**
 * @param {!Array<!proto.fm.internal.KeyLayoutBinding>} value
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
*/
proto.fm.internal.CharacterSaveEntry.prototype.setKeyLayoutList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 6, value);
};


/**
 * @param {!proto.fm.internal.KeyLayoutBinding=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.KeyLayoutBinding}
 */
proto.fm.internal.CharacterSaveEntry.prototype.addKeyLayout = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 6, opt_value, proto.fm.internal.KeyLayoutBinding, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.CharacterSaveEntry} returns this
 */
proto.fm.internal.CharacterSaveEntry.prototype.clearKeyLayoutList = function() {
  return this.setKeyLayoutList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.fm.internal.SaveCharactersRequest.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.SaveCharactersRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.SaveCharactersRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.SaveCharactersRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharactersRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    entriesList: jspb.Message.toObjectList(msg.getEntriesList(),
    proto.fm.internal.CharacterSaveEntry.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.SaveCharactersRequest}
 */
proto.fm.internal.SaveCharactersRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.SaveCharactersRequest;
  return proto.fm.internal.SaveCharactersRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.SaveCharactersRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.SaveCharactersRequest}
 */
proto.fm.internal.SaveCharactersRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.fm.internal.CharacterSaveEntry;
      reader.readMessage(value,proto.fm.internal.CharacterSaveEntry.deserializeBinaryFromReader);
      msg.addEntries(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.SaveCharactersRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.SaveCharactersRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.SaveCharactersRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharactersRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEntriesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.fm.internal.CharacterSaveEntry.serializeBinaryToWriter
    );
  }
};


/**
 * repeated CharacterSaveEntry entries = 1;
 * @return {!Array<!proto.fm.internal.CharacterSaveEntry>}
 */
proto.fm.internal.SaveCharactersRequest.prototype.getEntriesList = function() {
  return /** @type{!Array<!proto.fm.internal.CharacterSaveEntry>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.CharacterSaveEntry, 1));
};


/**
 * @param {!Array<!proto.fm.internal.CharacterSaveEntry>} value
 * @return {!proto.fm.internal.SaveCharactersRequest} returns this
*/
proto.fm.internal.SaveCharactersRequest.prototype.setEntriesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.fm.internal.CharacterSaveEntry=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.CharacterSaveEntry}
 */
proto.fm.internal.SaveCharactersRequest.prototype.addEntries = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.fm.internal.CharacterSaveEntry, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.SaveCharactersRequest} returns this
 */
proto.fm.internal.SaveCharactersRequest.prototype.clearEntriesList = function() {
  return this.setEntriesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.SaveCharactersReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.SaveCharactersReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.SaveCharactersReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharactersReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.SaveCharactersReply}
 */
proto.fm.internal.SaveCharactersReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.SaveCharactersReply;
  return proto.fm.internal.SaveCharactersReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.SaveCharactersReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.SaveCharactersReply}
 */
proto.fm.internal.SaveCharactersReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.SaveCharactersReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.SaveCharactersReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.SaveCharactersReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SaveCharactersReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.SaveCharactersReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.SaveCharactersReply} returns this
 */
proto.fm.internal.SaveCharactersReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.InventoryPersisted.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.InventoryPersisted.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.InventoryPersisted} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.InventoryPersisted.toObject = function(includeInstance, msg) {
  var f, obj = {
    uniqueId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    ownerId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    itemId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    slot: jspb.Message.getFieldWithDefault(msg, 4, 0),
    count: jspb.Message.getFieldWithDefault(msg, 5, 0),
    expirationUnixMs: jspb.Message.getFieldWithDefault(msg, 6, 0),
    enchantChance: jspb.Message.getFieldWithDefault(msg, 7, 0),
    flag: jspb.Message.getFieldWithDefault(msg, 8, 0),
    skillBonus: jspb.Message.getFieldWithDefault(msg, 9, 0),
    ownerName: jspb.Message.getFieldWithDefault(msg, 10, ""),
    inventoryType: jspb.Message.getFieldWithDefault(msg, 11, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.InventoryPersisted}
 */
proto.fm.internal.InventoryPersisted.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.InventoryPersisted;
  return proto.fm.internal.InventoryPersisted.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.InventoryPersisted} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.InventoryPersisted}
 */
proto.fm.internal.InventoryPersisted.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setUniqueId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setOwnerId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setItemId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setSlot(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCount(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setExpirationUnixMs(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setEnchantChance(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setFlag(value);
      break;
    case 9:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSkillBonus(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setOwnerName(value);
      break;
    case 11:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setInventoryType(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.InventoryPersisted.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.InventoryPersisted.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.InventoryPersisted} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.InventoryPersisted.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {number} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getOwnerId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getItemId();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getSlot();
  if (f !== 0) {
    writer.writeInt32(
      4,
      f
    );
  }
  f = message.getCount();
  if (f !== 0) {
    writer.writeUint32(
      5,
      f
    );
  }
  f = message.getExpirationUnixMs();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getEnchantChance();
  if (f !== 0) {
    writer.writeUint32(
      7,
      f
    );
  }
  f = message.getFlag();
  if (f !== 0) {
    writer.writeUint32(
      8,
      f
    );
  }
  f = message.getSkillBonus();
  if (f !== 0) {
    writer.writeUint32(
      9,
      f
    );
  }
  f = message.getOwnerName();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getInventoryType();
  if (f !== 0) {
    writer.writeUint32(
      11,
      f
    );
  }
};


/**
 * optional uint64 unique_id = 1;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getUniqueId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setUniqueId = function(value) {
  return jspb.Message.setField(this, 1, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.clearUniqueId = function() {
  return jspb.Message.setField(this, 1, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.InventoryPersisted.prototype.hasUniqueId = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional uint32 owner_id = 2;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getOwnerId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setOwnerId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint32 item_id = 3;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getItemId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setItemId = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional int32 slot = 4;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getSlot = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setSlot = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional uint32 count = 5;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setCount = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional int64 expiration_unix_ms = 6;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getExpirationUnixMs = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setExpirationUnixMs = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional uint32 enchant_chance = 7;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getEnchantChance = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setEnchantChance = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional uint32 flag = 8;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getFlag = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setFlag = function(value) {
  return jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * optional uint32 skill_bonus = 9;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getSkillBonus = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setSkillBonus = function(value) {
  return jspb.Message.setProto3IntField(this, 9, value);
};


/**
 * optional string owner_name = 10;
 * @return {string}
 */
proto.fm.internal.InventoryPersisted.prototype.getOwnerName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setOwnerName = function(value) {
  return jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional uint32 inventory_type = 11;
 * @return {number}
 */
proto.fm.internal.InventoryPersisted.prototype.getInventoryType = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InventoryPersisted} returns this
 */
proto.fm.internal.InventoryPersisted.prototype.setInventoryType = function(value) {
  return jspb.Message.setProto3IntField(this, 11, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.SkillPersisted.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.SkillPersisted.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.SkillPersisted} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SkillPersisted.toObject = function(includeInstance, msg) {
  var f, obj = {
    characterId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    skillId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    level: jspb.Message.getFieldWithDefault(msg, 3, 0),
    masterLevel: jspb.Message.getFieldWithDefault(msg, 4, 0),
    cooldownEndUnixMs: jspb.Message.getFieldWithDefault(msg, 5, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.SkillPersisted}
 */
proto.fm.internal.SkillPersisted.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.SkillPersisted;
  return proto.fm.internal.SkillPersisted.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.SkillPersisted} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.SkillPersisted}
 */
proto.fm.internal.SkillPersisted.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSkillId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setLevel(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setMasterLevel(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCooldownEndUnixMs(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.SkillPersisted.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.SkillPersisted.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.SkillPersisted} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.SkillPersisted.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getSkillId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getLevel();
  if (f !== 0) {
    writer.writeInt32(
      3,
      f
    );
  }
  f = message.getMasterLevel();
  if (f !== 0) {
    writer.writeInt32(
      4,
      f
    );
  }
  f = message.getCooldownEndUnixMs();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
};


/**
 * optional uint32 character_id = 1;
 * @return {number}
 */
proto.fm.internal.SkillPersisted.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.SkillPersisted} returns this
 */
proto.fm.internal.SkillPersisted.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 skill_id = 2;
 * @return {number}
 */
proto.fm.internal.SkillPersisted.prototype.getSkillId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.SkillPersisted} returns this
 */
proto.fm.internal.SkillPersisted.prototype.setSkillId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int32 level = 3;
 * @return {number}
 */
proto.fm.internal.SkillPersisted.prototype.getLevel = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.SkillPersisted} returns this
 */
proto.fm.internal.SkillPersisted.prototype.setLevel = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional int32 master_level = 4;
 * @return {number}
 */
proto.fm.internal.SkillPersisted.prototype.getMasterLevel = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.SkillPersisted} returns this
 */
proto.fm.internal.SkillPersisted.prototype.setMasterLevel = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional int64 cooldown_end_unix_ms = 5;
 * @return {number}
 */
proto.fm.internal.SkillPersisted.prototype.getCooldownEndUnixMs = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.SkillPersisted} returns this
 */
proto.fm.internal.SkillPersisted.prototype.setCooldownEndUnixMs = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.LoginAccountRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.LoginAccountRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.LoginAccountRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LoginAccountRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    loginId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    password: jspb.Message.getFieldWithDefault(msg, 2, ""),
    macAddress: jspb.Message.getFieldWithDefault(msg, 3, ""),
    ipAddress: jspb.Message.getFieldWithDefault(msg, 4, ""),
    initialRole: jspb.Message.getFieldWithDefault(msg, 5, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.LoginAccountRequest}
 */
proto.fm.internal.LoginAccountRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.LoginAccountRequest;
  return proto.fm.internal.LoginAccountRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.LoginAccountRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.LoginAccountRequest}
 */
proto.fm.internal.LoginAccountRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setLoginId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPassword(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setMacAddress(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setIpAddress(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setInitialRole(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.LoginAccountRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.LoginAccountRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.LoginAccountRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LoginAccountRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getLoginId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPassword();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getMacAddress();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getIpAddress();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getInitialRole();
  if (f !== 0) {
    writer.writeUint32(
      5,
      f
    );
  }
};


/**
 * optional string login_id = 1;
 * @return {string}
 */
proto.fm.internal.LoginAccountRequest.prototype.getLoginId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.LoginAccountRequest} returns this
 */
proto.fm.internal.LoginAccountRequest.prototype.setLoginId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string password = 2;
 * @return {string}
 */
proto.fm.internal.LoginAccountRequest.prototype.getPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.LoginAccountRequest} returns this
 */
proto.fm.internal.LoginAccountRequest.prototype.setPassword = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string mac_address = 3;
 * @return {string}
 */
proto.fm.internal.LoginAccountRequest.prototype.getMacAddress = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.LoginAccountRequest} returns this
 */
proto.fm.internal.LoginAccountRequest.prototype.setMacAddress = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string ip_address = 4;
 * @return {string}
 */
proto.fm.internal.LoginAccountRequest.prototype.getIpAddress = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.LoginAccountRequest} returns this
 */
proto.fm.internal.LoginAccountRequest.prototype.setIpAddress = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional uint32 initial_role = 5;
 * @return {number}
 */
proto.fm.internal.LoginAccountRequest.prototype.getInitialRole = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LoginAccountRequest} returns this
 */
proto.fm.internal.LoginAccountRequest.prototype.setInitialRole = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.LoginAccountReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.LoginAccountReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.LoginAccountReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LoginAccountReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    status: jspb.Message.getFieldWithDefault(msg, 1, 0),
    accountId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    gender: jspb.Message.getFieldWithDefault(msg, 3, 0),
    isChatBlocked: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
    chatBlockedUntilUnixMs: jspb.Message.getFieldWithDefault(msg, 5, 0),
    characterSlotCount: jspb.Message.getFieldWithDefault(msg, 6, 0),
    role: jspb.Message.getFieldWithDefault(msg, 7, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.LoginAccountReply}
 */
proto.fm.internal.LoginAccountReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.LoginAccountReply;
  return proto.fm.internal.LoginAccountReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.LoginAccountReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.LoginAccountReply}
 */
proto.fm.internal.LoginAccountReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.fm.internal.LoginAccountReply.Status} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setGender(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsChatBlocked(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setChatBlockedUntilUnixMs(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterSlotCount(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setRole(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.LoginAccountReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.LoginAccountReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.LoginAccountReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LoginAccountReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getGender();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getIsChatBlocked();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getChatBlockedUntilUnixMs();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
  f = message.getCharacterSlotCount();
  if (f !== 0) {
    writer.writeUint32(
      6,
      f
    );
  }
  f = message.getRole();
  if (f !== 0) {
    writer.writeUint32(
      7,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.fm.internal.LoginAccountReply.Status = {
  SUCCESS: 0,
  WRONG_PASSWORD: 1,
  BANNED: 2,
  REGISTERED: 3,
  ALREADY_LOGGED_IN: 4
};

/**
 * optional Status status = 1;
 * @return {!proto.fm.internal.LoginAccountReply.Status}
 */
proto.fm.internal.LoginAccountReply.prototype.getStatus = function() {
  return /** @type {!proto.fm.internal.LoginAccountReply.Status} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.fm.internal.LoginAccountReply.Status} value
 * @return {!proto.fm.internal.LoginAccountReply} returns this
 */
proto.fm.internal.LoginAccountReply.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional uint32 account_id = 2;
 * @return {number}
 */
proto.fm.internal.LoginAccountReply.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LoginAccountReply} returns this
 */
proto.fm.internal.LoginAccountReply.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint32 gender = 3;
 * @return {number}
 */
proto.fm.internal.LoginAccountReply.prototype.getGender = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LoginAccountReply} returns this
 */
proto.fm.internal.LoginAccountReply.prototype.setGender = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional bool is_chat_blocked = 4;
 * @return {boolean}
 */
proto.fm.internal.LoginAccountReply.prototype.getIsChatBlocked = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.LoginAccountReply} returns this
 */
proto.fm.internal.LoginAccountReply.prototype.setIsChatBlocked = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * optional int64 chat_blocked_until_unix_ms = 5;
 * @return {number}
 */
proto.fm.internal.LoginAccountReply.prototype.getChatBlockedUntilUnixMs = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LoginAccountReply} returns this
 */
proto.fm.internal.LoginAccountReply.prototype.setChatBlockedUntilUnixMs = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional uint32 character_slot_count = 6;
 * @return {number}
 */
proto.fm.internal.LoginAccountReply.prototype.getCharacterSlotCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LoginAccountReply} returns this
 */
proto.fm.internal.LoginAccountReply.prototype.setCharacterSlotCount = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional uint32 role = 7;
 * @return {number}
 */
proto.fm.internal.LoginAccountReply.prototype.getRole = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LoginAccountReply} returns this
 */
proto.fm.internal.LoginAccountReply.prototype.setRole = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CharacterOverview.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CharacterOverview.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CharacterOverview} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CharacterOverview.toObject = function(includeInstance, msg) {
  var f, obj = {
    characterId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    gender: jspb.Message.getFieldWithDefault(msg, 3, 0),
    skinColor: jspb.Message.getFieldWithDefault(msg, 4, 0),
    face: jspb.Message.getFieldWithDefault(msg, 5, 0),
    hair: jspb.Message.getFieldWithDefault(msg, 6, 0),
    level: jspb.Message.getFieldWithDefault(msg, 7, 0),
    classId: jspb.Message.getFieldWithDefault(msg, 8, 0),
    mapId: jspb.Message.getFieldWithDefault(msg, 9, 0),
    spawnPoint: jspb.Message.getFieldWithDefault(msg, 10, 0),
    accountId: jspb.Message.getFieldWithDefault(msg, 11, 0),
    worldId: jspb.Message.getFieldWithDefault(msg, 12, 0),
    rank: jspb.Message.getFieldWithDefault(msg, 13, 0),
    rankDiff: jspb.Message.getFieldWithDefault(msg, 14, 0),
    classRank: jspb.Message.getFieldWithDefault(msg, 15, 0),
    classRankDiff: jspb.Message.getFieldWithDefault(msg, 16, 0),
    baseLooksMap: (f = msg.getBaseLooksMap()) ? f.toObject(includeInstance, undefined) : [],
    overlaysMap: (f = msg.getOverlaysMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CharacterOverview}
 */
proto.fm.internal.CharacterOverview.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CharacterOverview;
  return proto.fm.internal.CharacterOverview.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CharacterOverview} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CharacterOverview}
 */
proto.fm.internal.CharacterOverview.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setGender(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSkinColor(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setFace(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setHair(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setLevel(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setClassId(value);
      break;
    case 9:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMapId(value);
      break;
    case 10:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSpawnPoint(value);
      break;
    case 11:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 12:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 13:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setRank(value);
      break;
    case 14:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setRankDiff(value);
      break;
    case 15:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setClassRank(value);
      break;
    case 16:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setClassRankDiff(value);
      break;
    case 17:
      var value = msg.getBaseLooksMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readInt32, jspb.BinaryReader.prototype.readUint32, null, 0, 0);
         });
      break;
    case 18:
      var value = msg.getOverlaysMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readInt32, jspb.BinaryReader.prototype.readUint32, null, 0, 0);
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CharacterOverview.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CharacterOverview.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CharacterOverview} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CharacterOverview.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getGender();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getSkinColor();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
  f = message.getFace();
  if (f !== 0) {
    writer.writeUint32(
      5,
      f
    );
  }
  f = message.getHair();
  if (f !== 0) {
    writer.writeUint32(
      6,
      f
    );
  }
  f = message.getLevel();
  if (f !== 0) {
    writer.writeUint32(
      7,
      f
    );
  }
  f = message.getClassId();
  if (f !== 0) {
    writer.writeUint32(
      8,
      f
    );
  }
  f = message.getMapId();
  if (f !== 0) {
    writer.writeUint32(
      9,
      f
    );
  }
  f = message.getSpawnPoint();
  if (f !== 0) {
    writer.writeUint32(
      10,
      f
    );
  }
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      11,
      f
    );
  }
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      12,
      f
    );
  }
  f = message.getRank();
  if (f !== 0) {
    writer.writeInt32(
      13,
      f
    );
  }
  f = message.getRankDiff();
  if (f !== 0) {
    writer.writeInt32(
      14,
      f
    );
  }
  f = message.getClassRank();
  if (f !== 0) {
    writer.writeInt32(
      15,
      f
    );
  }
  f = message.getClassRankDiff();
  if (f !== 0) {
    writer.writeInt32(
      16,
      f
    );
  }
  f = message.getBaseLooksMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(17, writer, jspb.BinaryWriter.prototype.writeInt32, jspb.BinaryWriter.prototype.writeUint32);
  }
  f = message.getOverlaysMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(18, writer, jspb.BinaryWriter.prototype.writeInt32, jspb.BinaryWriter.prototype.writeUint32);
  }
};


/**
 * optional uint32 character_id = 1;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.fm.internal.CharacterOverview.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional uint32 gender = 3;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getGender = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setGender = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional uint32 skin_color = 4;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getSkinColor = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setSkinColor = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional uint32 face = 5;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getFace = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setFace = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional uint32 hair = 6;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getHair = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setHair = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional uint32 level = 7;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getLevel = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setLevel = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional uint32 class_id = 8;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getClassId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setClassId = function(value) {
  return jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * optional uint32 map_id = 9;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getMapId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setMapId = function(value) {
  return jspb.Message.setProto3IntField(this, 9, value);
};


/**
 * optional uint32 spawn_point = 10;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getSpawnPoint = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 10, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setSpawnPoint = function(value) {
  return jspb.Message.setProto3IntField(this, 10, value);
};


/**
 * optional uint32 account_id = 11;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 11, value);
};


/**
 * optional uint32 world_id = 12;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 12, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 12, value);
};


/**
 * optional int32 rank = 13;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getRank = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 13, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setRank = function(value) {
  return jspb.Message.setProto3IntField(this, 13, value);
};


/**
 * optional int32 rank_diff = 14;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getRankDiff = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 14, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setRankDiff = function(value) {
  return jspb.Message.setProto3IntField(this, 14, value);
};


/**
 * optional int32 class_rank = 15;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getClassRank = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 15, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setClassRank = function(value) {
  return jspb.Message.setProto3IntField(this, 15, value);
};


/**
 * optional int32 class_rank_diff = 16;
 * @return {number}
 */
proto.fm.internal.CharacterOverview.prototype.getClassRankDiff = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 16, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.setClassRankDiff = function(value) {
  return jspb.Message.setProto3IntField(this, 16, value);
};


/**
 * map<int32, uint32> base_looks = 17;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<number,number>}
 */
proto.fm.internal.CharacterOverview.prototype.getBaseLooksMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<number,number>} */ (
      jspb.Message.getMapField(this, 17, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.clearBaseLooksMap = function() {
  this.getBaseLooksMap().clear();
  return this;};


/**
 * map<int32, uint32> overlays = 18;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<number,number>}
 */
proto.fm.internal.CharacterOverview.prototype.getOverlaysMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<number,number>} */ (
      jspb.Message.getMapField(this, 18, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.fm.internal.CharacterOverview} returns this
 */
proto.fm.internal.CharacterOverview.prototype.clearOverlaysMap = function() {
  this.getOverlaysMap().clear();
  return this;};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.GetCharacterListRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.GetCharacterListRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.GetCharacterListRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetCharacterListRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    accountId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    worldId: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.GetCharacterListRequest}
 */
proto.fm.internal.GetCharacterListRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.GetCharacterListRequest;
  return proto.fm.internal.GetCharacterListRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.GetCharacterListRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.GetCharacterListRequest}
 */
proto.fm.internal.GetCharacterListRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.GetCharacterListRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.GetCharacterListRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.GetCharacterListRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetCharacterListRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
};


/**
 * optional uint32 account_id = 1;
 * @return {number}
 */
proto.fm.internal.GetCharacterListRequest.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.GetCharacterListRequest} returns this
 */
proto.fm.internal.GetCharacterListRequest.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 world_id = 2;
 * @return {number}
 */
proto.fm.internal.GetCharacterListRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.GetCharacterListRequest} returns this
 */
proto.fm.internal.GetCharacterListRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.fm.internal.GetCharacterListReply.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.GetCharacterListReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.GetCharacterListReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.GetCharacterListReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetCharacterListReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    charactersList: jspb.Message.toObjectList(msg.getCharactersList(),
    proto.fm.internal.CharacterOverview.toObject, includeInstance),
    slotCount: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.GetCharacterListReply}
 */
proto.fm.internal.GetCharacterListReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.GetCharacterListReply;
  return proto.fm.internal.GetCharacterListReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.GetCharacterListReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.GetCharacterListReply}
 */
proto.fm.internal.GetCharacterListReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.fm.internal.CharacterOverview;
      reader.readMessage(value,proto.fm.internal.CharacterOverview.deserializeBinaryFromReader);
      msg.addCharacters(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSlotCount(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.GetCharacterListReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.GetCharacterListReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.GetCharacterListReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetCharacterListReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCharactersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.fm.internal.CharacterOverview.serializeBinaryToWriter
    );
  }
  f = message.getSlotCount();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
};


/**
 * repeated CharacterOverview characters = 1;
 * @return {!Array<!proto.fm.internal.CharacterOverview>}
 */
proto.fm.internal.GetCharacterListReply.prototype.getCharactersList = function() {
  return /** @type{!Array<!proto.fm.internal.CharacterOverview>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.CharacterOverview, 1));
};


/**
 * @param {!Array<!proto.fm.internal.CharacterOverview>} value
 * @return {!proto.fm.internal.GetCharacterListReply} returns this
*/
proto.fm.internal.GetCharacterListReply.prototype.setCharactersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.fm.internal.CharacterOverview=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.CharacterOverview}
 */
proto.fm.internal.GetCharacterListReply.prototype.addCharacters = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.fm.internal.CharacterOverview, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.GetCharacterListReply} returns this
 */
proto.fm.internal.GetCharacterListReply.prototype.clearCharactersList = function() {
  return this.setCharactersList([]);
};


/**
 * optional uint32 slot_count = 2;
 * @return {number}
 */
proto.fm.internal.GetCharacterListReply.prototype.getSlotCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.GetCharacterListReply} returns this
 */
proto.fm.internal.GetCharacterListReply.prototype.setSlotCount = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CheckCharacterNameRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CheckCharacterNameRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CheckCharacterNameRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CheckCharacterNameRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CheckCharacterNameRequest}
 */
proto.fm.internal.CheckCharacterNameRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CheckCharacterNameRequest;
  return proto.fm.internal.CheckCharacterNameRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CheckCharacterNameRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CheckCharacterNameRequest}
 */
proto.fm.internal.CheckCharacterNameRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CheckCharacterNameRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CheckCharacterNameRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CheckCharacterNameRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CheckCharacterNameRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.fm.internal.CheckCharacterNameRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.CheckCharacterNameRequest} returns this
 */
proto.fm.internal.CheckCharacterNameRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CheckCharacterNameReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CheckCharacterNameReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CheckCharacterNameReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CheckCharacterNameReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    exists: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CheckCharacterNameReply}
 */
proto.fm.internal.CheckCharacterNameReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CheckCharacterNameReply;
  return proto.fm.internal.CheckCharacterNameReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CheckCharacterNameReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CheckCharacterNameReply}
 */
proto.fm.internal.CheckCharacterNameReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setExists(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CheckCharacterNameReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CheckCharacterNameReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CheckCharacterNameReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CheckCharacterNameReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getExists();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool exists = 1;
 * @return {boolean}
 */
proto.fm.internal.CheckCharacterNameReply.prototype.getExists = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.CheckCharacterNameReply} returns this
 */
proto.fm.internal.CheckCharacterNameReply.prototype.setExists = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CreateCharacterRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CreateCharacterRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CreateCharacterRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreateCharacterRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    accountId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    worldId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    name: jspb.Message.getFieldWithDefault(msg, 3, ""),
    face: jspb.Message.getFieldWithDefault(msg, 4, 0),
    hair: jspb.Message.getFieldWithDefault(msg, 5, 0),
    skinColor: jspb.Message.getFieldWithDefault(msg, 6, 0),
    topItemId: jspb.Message.getFieldWithDefault(msg, 7, 0),
    bottomItemId: jspb.Message.getFieldWithDefault(msg, 8, 0),
    shoesItemId: jspb.Message.getFieldWithDefault(msg, 9, 0),
    weaponItemId: jspb.Message.getFieldWithDefault(msg, 10, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CreateCharacterRequest}
 */
proto.fm.internal.CreateCharacterRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CreateCharacterRequest;
  return proto.fm.internal.CreateCharacterRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CreateCharacterRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CreateCharacterRequest}
 */
proto.fm.internal.CreateCharacterRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setFace(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setHair(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setSkinColor(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setTopItemId(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setBottomItemId(value);
      break;
    case 9:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setShoesItemId(value);
      break;
    case 10:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWeaponItemId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CreateCharacterRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CreateCharacterRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CreateCharacterRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreateCharacterRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getFace();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
  f = message.getHair();
  if (f !== 0) {
    writer.writeUint32(
      5,
      f
    );
  }
  f = message.getSkinColor();
  if (f !== 0) {
    writer.writeUint32(
      6,
      f
    );
  }
  f = message.getTopItemId();
  if (f !== 0) {
    writer.writeUint32(
      7,
      f
    );
  }
  f = message.getBottomItemId();
  if (f !== 0) {
    writer.writeUint32(
      8,
      f
    );
  }
  f = message.getShoesItemId();
  if (f !== 0) {
    writer.writeUint32(
      9,
      f
    );
  }
  f = message.getWeaponItemId();
  if (f !== 0) {
    writer.writeUint32(
      10,
      f
    );
  }
};


/**
 * optional uint32 account_id = 1;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 world_id = 2;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional uint32 face = 4;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getFace = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setFace = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional uint32 hair = 5;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getHair = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setHair = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional uint32 skin_color = 6;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getSkinColor = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setSkinColor = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional uint32 top_item_id = 7;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getTopItemId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setTopItemId = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional uint32 bottom_item_id = 8;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getBottomItemId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setBottomItemId = function(value) {
  return jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * optional uint32 shoes_item_id = 9;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getShoesItemId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setShoesItemId = function(value) {
  return jspb.Message.setProto3IntField(this, 9, value);
};


/**
 * optional uint32 weapon_item_id = 10;
 * @return {number}
 */
proto.fm.internal.CreateCharacterRequest.prototype.getWeaponItemId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 10, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreateCharacterRequest} returns this
 */
proto.fm.internal.CreateCharacterRequest.prototype.setWeaponItemId = function(value) {
  return jspb.Message.setProto3IntField(this, 10, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CreateCharacterReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CreateCharacterReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CreateCharacterReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreateCharacterReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    success: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorMsg: jspb.Message.getFieldWithDefault(msg, 2, ""),
    character: (f = msg.getCharacter()) && proto.fm.internal.CharacterOverview.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CreateCharacterReply}
 */
proto.fm.internal.CreateCharacterReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CreateCharacterReply;
  return proto.fm.internal.CreateCharacterReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CreateCharacterReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CreateCharacterReply}
 */
proto.fm.internal.CreateCharacterReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setSuccess(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setErrorMsg(value);
      break;
    case 3:
      var value = new proto.fm.internal.CharacterOverview;
      reader.readMessage(value,proto.fm.internal.CharacterOverview.deserializeBinaryFromReader);
      msg.setCharacter(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CreateCharacterReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CreateCharacterReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CreateCharacterReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreateCharacterReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSuccess();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorMsg();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getCharacter();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.fm.internal.CharacterOverview.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool success = 1;
 * @return {boolean}
 */
proto.fm.internal.CreateCharacterReply.prototype.getSuccess = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.CreateCharacterReply} returns this
 */
proto.fm.internal.CreateCharacterReply.prototype.setSuccess = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional string error_msg = 2;
 * @return {string}
 */
proto.fm.internal.CreateCharacterReply.prototype.getErrorMsg = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.CreateCharacterReply} returns this
 */
proto.fm.internal.CreateCharacterReply.prototype.setErrorMsg = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional CharacterOverview character = 3;
 * @return {?proto.fm.internal.CharacterOverview}
 */
proto.fm.internal.CreateCharacterReply.prototype.getCharacter = function() {
  return /** @type{?proto.fm.internal.CharacterOverview} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.CharacterOverview, 3));
};


/**
 * @param {?proto.fm.internal.CharacterOverview|undefined} value
 * @return {!proto.fm.internal.CreateCharacterReply} returns this
*/
proto.fm.internal.CreateCharacterReply.prototype.setCharacter = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.CreateCharacterReply} returns this
 */
proto.fm.internal.CreateCharacterReply.prototype.clearCharacter = function() {
  return this.setCharacter(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.CreateCharacterReply.prototype.hasCharacter = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.DeleteCharacterRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.DeleteCharacterRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.DeleteCharacterRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DeleteCharacterRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    accountId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    characterId: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.DeleteCharacterRequest}
 */
proto.fm.internal.DeleteCharacterRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.DeleteCharacterRequest;
  return proto.fm.internal.DeleteCharacterRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.DeleteCharacterRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.DeleteCharacterRequest}
 */
proto.fm.internal.DeleteCharacterRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.DeleteCharacterRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.DeleteCharacterRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.DeleteCharacterRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DeleteCharacterRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
};


/**
 * optional uint32 account_id = 1;
 * @return {number}
 */
proto.fm.internal.DeleteCharacterRequest.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.DeleteCharacterRequest} returns this
 */
proto.fm.internal.DeleteCharacterRequest.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 character_id = 2;
 * @return {number}
 */
proto.fm.internal.DeleteCharacterRequest.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.DeleteCharacterRequest} returns this
 */
proto.fm.internal.DeleteCharacterRequest.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.DeleteCharacterReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.DeleteCharacterReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.DeleteCharacterReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DeleteCharacterReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    success: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.DeleteCharacterReply}
 */
proto.fm.internal.DeleteCharacterReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.DeleteCharacterReply;
  return proto.fm.internal.DeleteCharacterReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.DeleteCharacterReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.DeleteCharacterReply}
 */
proto.fm.internal.DeleteCharacterReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setSuccess(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.DeleteCharacterReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.DeleteCharacterReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.DeleteCharacterReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DeleteCharacterReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSuccess();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool success = 1;
 * @return {boolean}
 */
proto.fm.internal.DeleteCharacterReply.prototype.getSuccess = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.DeleteCharacterReply} returns this
 */
proto.fm.internal.DeleteCharacterReply.prototype.setSuccess = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.RefreshSessionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.RefreshSessionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.RefreshSessionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.RefreshSessionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    accountId: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.RefreshSessionRequest}
 */
proto.fm.internal.RefreshSessionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.RefreshSessionRequest;
  return proto.fm.internal.RefreshSessionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.RefreshSessionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.RefreshSessionRequest}
 */
proto.fm.internal.RefreshSessionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.RefreshSessionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.RefreshSessionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.RefreshSessionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.RefreshSessionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.RefreshSessionRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.RefreshSessionRequest} returns this
 */
proto.fm.internal.RefreshSessionRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 account_id = 2;
 * @return {number}
 */
proto.fm.internal.RefreshSessionRequest.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.RefreshSessionRequest} returns this
 */
proto.fm.internal.RefreshSessionRequest.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.RefreshSessionReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.RefreshSessionReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.RefreshSessionReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.RefreshSessionReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.RefreshSessionReply}
 */
proto.fm.internal.RefreshSessionReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.RefreshSessionReply;
  return proto.fm.internal.RefreshSessionReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.RefreshSessionReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.RefreshSessionReply}
 */
proto.fm.internal.RefreshSessionReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.SessionErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.RefreshSessionReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.RefreshSessionReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.RefreshSessionReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.RefreshSessionReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.RefreshSessionReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.RefreshSessionReply} returns this
 */
proto.fm.internal.RefreshSessionReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional SessionErrorCode error_code = 2;
 * @return {!proto.fm.internal.SessionErrorCode}
 */
proto.fm.internal.RefreshSessionReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.SessionErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.SessionErrorCode} value
 * @return {!proto.fm.internal.RefreshSessionReply} returns this
 */
proto.fm.internal.RefreshSessionReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.LogoutSessionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.LogoutSessionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.LogoutSessionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LogoutSessionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    accountId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    disconnectSource: jspb.Message.getFieldWithDefault(msg, 3, 0),
    transferDisconnect: jspb.Message.getBooleanFieldWithDefault(msg, 4, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.LogoutSessionRequest}
 */
proto.fm.internal.LogoutSessionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.LogoutSessionRequest;
  return proto.fm.internal.LogoutSessionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.LogoutSessionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.LogoutSessionRequest}
 */
proto.fm.internal.LogoutSessionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAccountId(value);
      break;
    case 3:
      var value = /** @type {!proto.fm.internal.SessionDisconnectSource} */ (reader.readEnum());
      msg.setDisconnectSource(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setTransferDisconnect(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.LogoutSessionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.LogoutSessionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.LogoutSessionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LogoutSessionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getAccountId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getDisconnectSource();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getTransferDisconnect();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.LogoutSessionRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LogoutSessionRequest} returns this
 */
proto.fm.internal.LogoutSessionRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 account_id = 2;
 * @return {number}
 */
proto.fm.internal.LogoutSessionRequest.prototype.getAccountId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LogoutSessionRequest} returns this
 */
proto.fm.internal.LogoutSessionRequest.prototype.setAccountId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional SessionDisconnectSource disconnect_source = 3;
 * @return {!proto.fm.internal.SessionDisconnectSource}
 */
proto.fm.internal.LogoutSessionRequest.prototype.getDisconnectSource = function() {
  return /** @type {!proto.fm.internal.SessionDisconnectSource} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.fm.internal.SessionDisconnectSource} value
 * @return {!proto.fm.internal.LogoutSessionRequest} returns this
 */
proto.fm.internal.LogoutSessionRequest.prototype.setDisconnectSource = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional bool transfer_disconnect = 4;
 * @return {boolean}
 */
proto.fm.internal.LogoutSessionRequest.prototype.getTransferDisconnect = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.LogoutSessionRequest} returns this
 */
proto.fm.internal.LogoutSessionRequest.prototype.setTransferDisconnect = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.LogoutSessionReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.LogoutSessionReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.LogoutSessionReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LogoutSessionReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.LogoutSessionReply}
 */
proto.fm.internal.LogoutSessionReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.LogoutSessionReply;
  return proto.fm.internal.LogoutSessionReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.LogoutSessionReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.LogoutSessionReply}
 */
proto.fm.internal.LogoutSessionReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.LogoutSessionReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.LogoutSessionReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.LogoutSessionReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LogoutSessionReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.LogoutSessionReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.LogoutSessionReply} returns this
 */
proto.fm.internal.LogoutSessionReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.PartyDoor.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.PartyDoor.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.PartyDoor} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PartyDoor.toObject = function(includeInstance, msg) {
  var f, obj = {
    town: jspb.Message.getFieldWithDefault(msg, 1, 0),
    target: jspb.Message.getFieldWithDefault(msg, 2, 0),
    x: jspb.Message.getFieldWithDefault(msg, 3, 0),
    y: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.PartyDoor}
 */
proto.fm.internal.PartyDoor.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.PartyDoor;
  return proto.fm.internal.PartyDoor.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.PartyDoor} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.PartyDoor}
 */
proto.fm.internal.PartyDoor.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setTown(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setTarget(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setX(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setY(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.PartyDoor.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.PartyDoor.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.PartyDoor} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PartyDoor.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTown();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getTarget();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getX();
  if (f !== 0) {
    writer.writeInt32(
      3,
      f
    );
  }
  f = message.getY();
  if (f !== 0) {
    writer.writeInt32(
      4,
      f
    );
  }
};


/**
 * optional uint32 town = 1;
 * @return {number}
 */
proto.fm.internal.PartyDoor.prototype.getTown = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyDoor} returns this
 */
proto.fm.internal.PartyDoor.prototype.setTown = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 target = 2;
 * @return {number}
 */
proto.fm.internal.PartyDoor.prototype.getTarget = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyDoor} returns this
 */
proto.fm.internal.PartyDoor.prototype.setTarget = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int32 x = 3;
 * @return {number}
 */
proto.fm.internal.PartyDoor.prototype.getX = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyDoor} returns this
 */
proto.fm.internal.PartyDoor.prototype.setX = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional int32 y = 4;
 * @return {number}
 */
proto.fm.internal.PartyDoor.prototype.getY = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyDoor} returns this
 */
proto.fm.internal.PartyDoor.prototype.setY = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.PartyMember.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.PartyMember.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.PartyMember} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PartyMember.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    characterId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    characterName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    level: jspb.Message.getFieldWithDefault(msg, 4, 0),
    classId: jspb.Message.getFieldWithDefault(msg, 5, 0),
    role: jspb.Message.getFieldWithDefault(msg, 6, ""),
    mapId: jspb.Message.getFieldWithDefault(msg, 7, 0),
    channelIndex: jspb.Message.getFieldWithDefault(msg, 8, 0),
    door: (f = msg.getDoor()) && proto.fm.internal.PartyDoor.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.PartyMember}
 */
proto.fm.internal.PartyMember.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.PartyMember;
  return proto.fm.internal.PartyMember.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.PartyMember} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.PartyMember}
 */
proto.fm.internal.PartyMember.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCharacterName(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setLevel(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setClassId(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setRole(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMapId(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setChannelIndex(value);
      break;
    case 9:
      var value = new proto.fm.internal.PartyDoor;
      reader.readMessage(value,proto.fm.internal.PartyDoor.deserializeBinaryFromReader);
      msg.setDoor(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.PartyMember.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.PartyMember.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.PartyMember} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.PartyMember.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getCharacterName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getLevel();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
  f = message.getClassId();
  if (f !== 0) {
    writer.writeUint32(
      5,
      f
    );
  }
  f = message.getRole();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getMapId();
  if (f !== 0) {
    writer.writeUint32(
      7,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 8));
  if (f != null) {
    writer.writeInt32(
      8,
      f
    );
  }
  f = message.getDoor();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      proto.fm.internal.PartyDoor.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.PartyMember.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 character_id = 2;
 * @return {number}
 */
proto.fm.internal.PartyMember.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string character_name = 3;
 * @return {string}
 */
proto.fm.internal.PartyMember.prototype.getCharacterName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setCharacterName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional uint32 level = 4;
 * @return {number}
 */
proto.fm.internal.PartyMember.prototype.getLevel = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setLevel = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional uint32 class_id = 5;
 * @return {number}
 */
proto.fm.internal.PartyMember.prototype.getClassId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setClassId = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional string role = 6;
 * @return {string}
 */
proto.fm.internal.PartyMember.prototype.getRole = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setRole = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional uint32 map_id = 7;
 * @return {number}
 */
proto.fm.internal.PartyMember.prototype.getMapId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setMapId = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional int32 channel_index = 8;
 * @return {number}
 */
proto.fm.internal.PartyMember.prototype.getChannelIndex = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.setChannelIndex = function(value) {
  return jspb.Message.setField(this, 8, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.clearChannelIndex = function() {
  return jspb.Message.setField(this, 8, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.PartyMember.prototype.hasChannelIndex = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional PartyDoor door = 9;
 * @return {?proto.fm.internal.PartyDoor}
 */
proto.fm.internal.PartyMember.prototype.getDoor = function() {
  return /** @type{?proto.fm.internal.PartyDoor} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.PartyDoor, 9));
};


/**
 * @param {?proto.fm.internal.PartyDoor|undefined} value
 * @return {!proto.fm.internal.PartyMember} returns this
*/
proto.fm.internal.PartyMember.prototype.setDoor = function(value) {
  return jspb.Message.setWrapperField(this, 9, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.PartyMember} returns this
 */
proto.fm.internal.PartyMember.prototype.clearDoor = function() {
  return this.setDoor(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.PartyMember.prototype.hasDoor = function() {
  return jspb.Message.getField(this, 9) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CreatePartyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CreatePartyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CreatePartyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreatePartyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    leader: (f = msg.getLeader()) && proto.fm.internal.PartyMember.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CreatePartyRequest}
 */
proto.fm.internal.CreatePartyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CreatePartyRequest;
  return proto.fm.internal.CreatePartyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CreatePartyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CreatePartyRequest}
 */
proto.fm.internal.CreatePartyRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = new proto.fm.internal.PartyMember;
      reader.readMessage(value,proto.fm.internal.PartyMember.deserializeBinaryFromReader);
      msg.setLeader(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CreatePartyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CreatePartyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CreatePartyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreatePartyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getLeader();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.fm.internal.PartyMember.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.CreatePartyRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreatePartyRequest} returns this
 */
proto.fm.internal.CreatePartyRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional PartyMember leader = 2;
 * @return {?proto.fm.internal.PartyMember}
 */
proto.fm.internal.CreatePartyRequest.prototype.getLeader = function() {
  return /** @type{?proto.fm.internal.PartyMember} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.PartyMember, 2));
};


/**
 * @param {?proto.fm.internal.PartyMember|undefined} value
 * @return {!proto.fm.internal.CreatePartyRequest} returns this
*/
proto.fm.internal.CreatePartyRequest.prototype.setLeader = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.CreatePartyRequest} returns this
 */
proto.fm.internal.CreatePartyRequest.prototype.clearLeader = function() {
  return this.setLeader(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.CreatePartyRequest.prototype.hasLeader = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.CreatePartyReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.CreatePartyReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.CreatePartyReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreatePartyReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    revision: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.CreatePartyReply}
 */
proto.fm.internal.CreatePartyReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.CreatePartyReply;
  return proto.fm.internal.CreatePartyReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.CreatePartyReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.CreatePartyReply}
 */
proto.fm.internal.CreatePartyReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setRevision(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.CreatePartyReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.CreatePartyReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.CreatePartyReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.CreatePartyReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getRevision();
  if (f !== 0) {
    writer.writeUint64(
      4,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.CreatePartyReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.CreatePartyReply} returns this
 */
proto.fm.internal.CreatePartyReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.CreatePartyReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.CreatePartyReply} returns this
 */
proto.fm.internal.CreatePartyReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional uint32 party_id = 3;
 * @return {number}
 */
proto.fm.internal.CreatePartyReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreatePartyReply} returns this
 */
proto.fm.internal.CreatePartyReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.CreatePartyReply} returns this
 */
proto.fm.internal.CreatePartyReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.CreatePartyReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional uint64 revision = 4;
 * @return {number}
 */
proto.fm.internal.CreatePartyReply.prototype.getRevision = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.CreatePartyReply} returns this
 */
proto.fm.internal.CreatePartyReply.prototype.setRevision = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.JoinPartyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.JoinPartyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.JoinPartyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.JoinPartyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    member: (f = msg.getMember()) && proto.fm.internal.PartyMember.toObject(includeInstance, f),
    skipInvitePendingCheck: jspb.Message.getBooleanFieldWithDefault(msg, 4, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.JoinPartyRequest}
 */
proto.fm.internal.JoinPartyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.JoinPartyRequest;
  return proto.fm.internal.JoinPartyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.JoinPartyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.JoinPartyRequest}
 */
proto.fm.internal.JoinPartyRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 3:
      var value = new proto.fm.internal.PartyMember;
      reader.readMessage(value,proto.fm.internal.PartyMember.deserializeBinaryFromReader);
      msg.setMember(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setSkipInvitePendingCheck(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.JoinPartyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.JoinPartyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.JoinPartyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.JoinPartyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getPartyId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getMember();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.fm.internal.PartyMember.serializeBinaryToWriter
    );
  }
  f = message.getSkipInvitePendingCheck();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.JoinPartyRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.JoinPartyRequest} returns this
 */
proto.fm.internal.JoinPartyRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 party_id = 2;
 * @return {number}
 */
proto.fm.internal.JoinPartyRequest.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.JoinPartyRequest} returns this
 */
proto.fm.internal.JoinPartyRequest.prototype.setPartyId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional PartyMember member = 3;
 * @return {?proto.fm.internal.PartyMember}
 */
proto.fm.internal.JoinPartyRequest.prototype.getMember = function() {
  return /** @type{?proto.fm.internal.PartyMember} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.PartyMember, 3));
};


/**
 * @param {?proto.fm.internal.PartyMember|undefined} value
 * @return {!proto.fm.internal.JoinPartyRequest} returns this
*/
proto.fm.internal.JoinPartyRequest.prototype.setMember = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.JoinPartyRequest} returns this
 */
proto.fm.internal.JoinPartyRequest.prototype.clearMember = function() {
  return this.setMember(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.JoinPartyRequest.prototype.hasMember = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional bool skip_invite_pending_check = 4;
 * @return {boolean}
 */
proto.fm.internal.JoinPartyRequest.prototype.getSkipInvitePendingCheck = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.JoinPartyRequest} returns this
 */
proto.fm.internal.JoinPartyRequest.prototype.setSkipInvitePendingCheck = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.JoinPartyReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.JoinPartyReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.JoinPartyReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.JoinPartyReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    revision: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.JoinPartyReply}
 */
proto.fm.internal.JoinPartyReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.JoinPartyReply;
  return proto.fm.internal.JoinPartyReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.JoinPartyReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.JoinPartyReply}
 */
proto.fm.internal.JoinPartyReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setRevision(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.JoinPartyReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.JoinPartyReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.JoinPartyReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.JoinPartyReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getRevision();
  if (f !== 0) {
    writer.writeUint64(
      4,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.JoinPartyReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.JoinPartyReply} returns this
 */
proto.fm.internal.JoinPartyReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.JoinPartyReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.JoinPartyReply} returns this
 */
proto.fm.internal.JoinPartyReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional uint32 party_id = 3;
 * @return {number}
 */
proto.fm.internal.JoinPartyReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.JoinPartyReply} returns this
 */
proto.fm.internal.JoinPartyReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.JoinPartyReply} returns this
 */
proto.fm.internal.JoinPartyReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.JoinPartyReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional uint64 revision = 4;
 * @return {number}
 */
proto.fm.internal.JoinPartyReply.prototype.getRevision = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.JoinPartyReply} returns this
 */
proto.fm.internal.JoinPartyReply.prototype.setRevision = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.LeavePartyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.LeavePartyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.LeavePartyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LeavePartyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    characterId: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.LeavePartyRequest}
 */
proto.fm.internal.LeavePartyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.LeavePartyRequest;
  return proto.fm.internal.LeavePartyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.LeavePartyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.LeavePartyRequest}
 */
proto.fm.internal.LeavePartyRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setCharacterId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.LeavePartyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.LeavePartyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.LeavePartyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LeavePartyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.LeavePartyRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LeavePartyRequest} returns this
 */
proto.fm.internal.LeavePartyRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 character_id = 2;
 * @return {number}
 */
proto.fm.internal.LeavePartyRequest.prototype.getCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LeavePartyRequest} returns this
 */
proto.fm.internal.LeavePartyRequest.prototype.setCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.LeavePartyReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.LeavePartyReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.LeavePartyReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LeavePartyReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    revision: jspb.Message.getFieldWithDefault(msg, 4, 0),
    disbanded: jspb.Message.getBooleanFieldWithDefault(msg, 5, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.LeavePartyReply}
 */
proto.fm.internal.LeavePartyReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.LeavePartyReply;
  return proto.fm.internal.LeavePartyReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.LeavePartyReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.LeavePartyReply}
 */
proto.fm.internal.LeavePartyReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setRevision(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisbanded(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.LeavePartyReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.LeavePartyReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.LeavePartyReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.LeavePartyReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getRevision();
  if (f !== 0) {
    writer.writeUint64(
      4,
      f
    );
  }
  f = message.getDisbanded();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.LeavePartyReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.LeavePartyReply} returns this
 */
proto.fm.internal.LeavePartyReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.LeavePartyReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.LeavePartyReply} returns this
 */
proto.fm.internal.LeavePartyReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional uint32 party_id = 3;
 * @return {number}
 */
proto.fm.internal.LeavePartyReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LeavePartyReply} returns this
 */
proto.fm.internal.LeavePartyReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.LeavePartyReply} returns this
 */
proto.fm.internal.LeavePartyReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.LeavePartyReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional uint64 revision = 4;
 * @return {number}
 */
proto.fm.internal.LeavePartyReply.prototype.getRevision = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.LeavePartyReply} returns this
 */
proto.fm.internal.LeavePartyReply.prototype.setRevision = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional bool disbanded = 5;
 * @return {boolean}
 */
proto.fm.internal.LeavePartyReply.prototype.getDisbanded = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 5, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.LeavePartyReply} returns this
 */
proto.fm.internal.LeavePartyReply.prototype.setDisbanded = function(value) {
  return jspb.Message.setProto3BooleanField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.ExpelPartyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.ExpelPartyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.ExpelPartyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ExpelPartyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    requesterCharacterId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    targetCharacterId: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.ExpelPartyRequest}
 */
proto.fm.internal.ExpelPartyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.ExpelPartyRequest;
  return proto.fm.internal.ExpelPartyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.ExpelPartyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.ExpelPartyRequest}
 */
proto.fm.internal.ExpelPartyRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setRequesterCharacterId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setTargetCharacterId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.ExpelPartyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.ExpelPartyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.ExpelPartyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ExpelPartyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getRequesterCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getTargetCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.ExpelPartyRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ExpelPartyRequest} returns this
 */
proto.fm.internal.ExpelPartyRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 requester_character_id = 2;
 * @return {number}
 */
proto.fm.internal.ExpelPartyRequest.prototype.getRequesterCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ExpelPartyRequest} returns this
 */
proto.fm.internal.ExpelPartyRequest.prototype.setRequesterCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint32 target_character_id = 3;
 * @return {number}
 */
proto.fm.internal.ExpelPartyRequest.prototype.getTargetCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ExpelPartyRequest} returns this
 */
proto.fm.internal.ExpelPartyRequest.prototype.setTargetCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.ExpelPartyReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.ExpelPartyReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.ExpelPartyReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ExpelPartyReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    revision: jspb.Message.getFieldWithDefault(msg, 4, 0),
    disbanded: jspb.Message.getBooleanFieldWithDefault(msg, 5, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.ExpelPartyReply}
 */
proto.fm.internal.ExpelPartyReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.ExpelPartyReply;
  return proto.fm.internal.ExpelPartyReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.ExpelPartyReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.ExpelPartyReply}
 */
proto.fm.internal.ExpelPartyReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setRevision(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisbanded(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.ExpelPartyReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.ExpelPartyReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.ExpelPartyReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ExpelPartyReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getRevision();
  if (f !== 0) {
    writer.writeUint64(
      4,
      f
    );
  }
  f = message.getDisbanded();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.ExpelPartyReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.ExpelPartyReply} returns this
 */
proto.fm.internal.ExpelPartyReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.ExpelPartyReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.ExpelPartyReply} returns this
 */
proto.fm.internal.ExpelPartyReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional uint32 party_id = 3;
 * @return {number}
 */
proto.fm.internal.ExpelPartyReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ExpelPartyReply} returns this
 */
proto.fm.internal.ExpelPartyReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.ExpelPartyReply} returns this
 */
proto.fm.internal.ExpelPartyReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.ExpelPartyReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional uint64 revision = 4;
 * @return {number}
 */
proto.fm.internal.ExpelPartyReply.prototype.getRevision = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ExpelPartyReply} returns this
 */
proto.fm.internal.ExpelPartyReply.prototype.setRevision = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional bool disbanded = 5;
 * @return {boolean}
 */
proto.fm.internal.ExpelPartyReply.prototype.getDisbanded = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 5, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.ExpelPartyReply} returns this
 */
proto.fm.internal.ExpelPartyReply.prototype.setDisbanded = function(value) {
  return jspb.Message.setProto3BooleanField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.ChangePartyLeaderRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.ChangePartyLeaderRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ChangePartyLeaderRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    requesterCharacterId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    newLeaderCharacterId: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.ChangePartyLeaderRequest}
 */
proto.fm.internal.ChangePartyLeaderRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.ChangePartyLeaderRequest;
  return proto.fm.internal.ChangePartyLeaderRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.ChangePartyLeaderRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.ChangePartyLeaderRequest}
 */
proto.fm.internal.ChangePartyLeaderRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setRequesterCharacterId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setNewLeaderCharacterId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.ChangePartyLeaderRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.ChangePartyLeaderRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ChangePartyLeaderRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getPartyId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getRequesterCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getNewLeaderCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChangePartyLeaderRequest} returns this
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 party_id = 2;
 * @return {number}
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChangePartyLeaderRequest} returns this
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.setPartyId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint32 requester_character_id = 3;
 * @return {number}
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.getRequesterCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChangePartyLeaderRequest} returns this
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.setRequesterCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional uint32 new_leader_character_id = 4;
 * @return {number}
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.getNewLeaderCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChangePartyLeaderRequest} returns this
 */
proto.fm.internal.ChangePartyLeaderRequest.prototype.setNewLeaderCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.ChangePartyLeaderReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.ChangePartyLeaderReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ChangePartyLeaderReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    revision: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.ChangePartyLeaderReply}
 */
proto.fm.internal.ChangePartyLeaderReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.ChangePartyLeaderReply;
  return proto.fm.internal.ChangePartyLeaderReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.ChangePartyLeaderReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.ChangePartyLeaderReply}
 */
proto.fm.internal.ChangePartyLeaderReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setRevision(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.ChangePartyLeaderReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.ChangePartyLeaderReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.ChangePartyLeaderReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getRevision();
  if (f !== 0) {
    writer.writeUint64(
      4,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.ChangePartyLeaderReply} returns this
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.ChangePartyLeaderReply} returns this
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional uint32 party_id = 3;
 * @return {number}
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChangePartyLeaderReply} returns this
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.ChangePartyLeaderReply} returns this
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional uint64 revision = 4;
 * @return {number}
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.getRevision = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.ChangePartyLeaderReply} returns this
 */
proto.fm.internal.ChangePartyLeaderReply.prototype.setRevision = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.fm.internal.Party.repeatedFields_ = [6];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.Party.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.Party.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.Party} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.Party.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    leaderCharacterId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    revision: jspb.Message.getFieldWithDefault(msg, 4, 0),
    state: jspb.Message.getFieldWithDefault(msg, 5, ""),
    membersList: jspb.Message.toObjectList(msg.getMembersList(),
    proto.fm.internal.PartyMember.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.Party}
 */
proto.fm.internal.Party.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.Party;
  return proto.fm.internal.Party.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.Party} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.Party}
 */
proto.fm.internal.Party.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setLeaderCharacterId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setRevision(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setState(value);
      break;
    case 6:
      var value = new proto.fm.internal.PartyMember;
      reader.readMessage(value,proto.fm.internal.PartyMember.deserializeBinaryFromReader);
      msg.addMembers(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.Party.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.Party.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.Party} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.Party.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getPartyId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getLeaderCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getRevision();
  if (f !== 0) {
    writer.writeUint64(
      4,
      f
    );
  }
  f = message.getState();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getMembersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      6,
      f,
      proto.fm.internal.PartyMember.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.Party.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.Party} returns this
 */
proto.fm.internal.Party.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 party_id = 2;
 * @return {number}
 */
proto.fm.internal.Party.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.Party} returns this
 */
proto.fm.internal.Party.prototype.setPartyId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint32 leader_character_id = 3;
 * @return {number}
 */
proto.fm.internal.Party.prototype.getLeaderCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.Party} returns this
 */
proto.fm.internal.Party.prototype.setLeaderCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional uint64 revision = 4;
 * @return {number}
 */
proto.fm.internal.Party.prototype.getRevision = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.Party} returns this
 */
proto.fm.internal.Party.prototype.setRevision = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional string state = 5;
 * @return {string}
 */
proto.fm.internal.Party.prototype.getState = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.Party} returns this
 */
proto.fm.internal.Party.prototype.setState = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * repeated PartyMember members = 6;
 * @return {!Array<!proto.fm.internal.PartyMember>}
 */
proto.fm.internal.Party.prototype.getMembersList = function() {
  return /** @type{!Array<!proto.fm.internal.PartyMember>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.fm.internal.PartyMember, 6));
};


/**
 * @param {!Array<!proto.fm.internal.PartyMember>} value
 * @return {!proto.fm.internal.Party} returns this
*/
proto.fm.internal.Party.prototype.setMembersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 6, value);
};


/**
 * @param {!proto.fm.internal.PartyMember=} opt_value
 * @param {number=} opt_index
 * @return {!proto.fm.internal.PartyMember}
 */
proto.fm.internal.Party.prototype.addMembers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 6, opt_value, proto.fm.internal.PartyMember, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.fm.internal.Party} returns this
 */
proto.fm.internal.Party.prototype.clearMembersList = function() {
  return this.setMembersList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.GetPartyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.GetPartyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.GetPartyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetPartyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.GetPartyRequest}
 */
proto.fm.internal.GetPartyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.GetPartyRequest;
  return proto.fm.internal.GetPartyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.GetPartyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.GetPartyRequest}
 */
proto.fm.internal.GetPartyRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.GetPartyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.GetPartyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.GetPartyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetPartyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getPartyId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.GetPartyRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.GetPartyRequest} returns this
 */
proto.fm.internal.GetPartyRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 party_id = 2;
 * @return {number}
 */
proto.fm.internal.GetPartyRequest.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.GetPartyRequest} returns this
 */
proto.fm.internal.GetPartyRequest.prototype.setPartyId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.GetPartyReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.GetPartyReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.GetPartyReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetPartyReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    found: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    party: (f = msg.getParty()) && proto.fm.internal.Party.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.GetPartyReply}
 */
proto.fm.internal.GetPartyReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.GetPartyReply;
  return proto.fm.internal.GetPartyReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.GetPartyReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.GetPartyReply}
 */
proto.fm.internal.GetPartyReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setFound(value);
      break;
    case 2:
      var value = new proto.fm.internal.Party;
      reader.readMessage(value,proto.fm.internal.Party.deserializeBinaryFromReader);
      msg.setParty(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.GetPartyReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.GetPartyReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.GetPartyReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.GetPartyReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFound();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getParty();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.fm.internal.Party.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool found = 1;
 * @return {boolean}
 */
proto.fm.internal.GetPartyReply.prototype.getFound = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.GetPartyReply} returns this
 */
proto.fm.internal.GetPartyReply.prototype.setFound = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional Party party = 2;
 * @return {?proto.fm.internal.Party}
 */
proto.fm.internal.GetPartyReply.prototype.getParty = function() {
  return /** @type{?proto.fm.internal.Party} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.Party, 2));
};


/**
 * @param {?proto.fm.internal.Party|undefined} value
 * @return {!proto.fm.internal.GetPartyReply} returns this
*/
proto.fm.internal.GetPartyReply.prototype.setParty = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.GetPartyReply} returns this
 */
proto.fm.internal.GetPartyReply.prototype.clearParty = function() {
  return this.setParty(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.GetPartyReply.prototype.hasParty = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.UpdatePartyMemberRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.UpdatePartyMemberRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.UpdatePartyMemberRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.UpdatePartyMemberRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    member: (f = msg.getMember()) && proto.fm.internal.PartyMember.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.UpdatePartyMemberRequest}
 */
proto.fm.internal.UpdatePartyMemberRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.UpdatePartyMemberRequest;
  return proto.fm.internal.UpdatePartyMemberRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.UpdatePartyMemberRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.UpdatePartyMemberRequest}
 */
proto.fm.internal.UpdatePartyMemberRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.fm.internal.PartyMember;
      reader.readMessage(value,proto.fm.internal.PartyMember.deserializeBinaryFromReader);
      msg.setMember(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.UpdatePartyMemberRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.UpdatePartyMemberRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.UpdatePartyMemberRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.UpdatePartyMemberRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMember();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.fm.internal.PartyMember.serializeBinaryToWriter
    );
  }
};


/**
 * optional PartyMember member = 1;
 * @return {?proto.fm.internal.PartyMember}
 */
proto.fm.internal.UpdatePartyMemberRequest.prototype.getMember = function() {
  return /** @type{?proto.fm.internal.PartyMember} */ (
    jspb.Message.getWrapperField(this, proto.fm.internal.PartyMember, 1));
};


/**
 * @param {?proto.fm.internal.PartyMember|undefined} value
 * @return {!proto.fm.internal.UpdatePartyMemberRequest} returns this
*/
proto.fm.internal.UpdatePartyMemberRequest.prototype.setMember = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.fm.internal.UpdatePartyMemberRequest} returns this
 */
proto.fm.internal.UpdatePartyMemberRequest.prototype.clearMember = function() {
  return this.setMember(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.UpdatePartyMemberRequest.prototype.hasMember = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.UpdatePartyMemberReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.UpdatePartyMemberReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.UpdatePartyMemberReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    revision: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.UpdatePartyMemberReply}
 */
proto.fm.internal.UpdatePartyMemberReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.UpdatePartyMemberReply;
  return proto.fm.internal.UpdatePartyMemberReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.UpdatePartyMemberReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.UpdatePartyMemberReply}
 */
proto.fm.internal.UpdatePartyMemberReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setRevision(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.UpdatePartyMemberReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.UpdatePartyMemberReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.UpdatePartyMemberReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getRevision();
  if (f !== 0) {
    writer.writeUint64(
      4,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.UpdatePartyMemberReply} returns this
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.UpdatePartyMemberReply} returns this
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional uint32 party_id = 3;
 * @return {number}
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.UpdatePartyMemberReply} returns this
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.UpdatePartyMemberReply} returns this
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional uint64 revision = 4;
 * @return {number}
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.getRevision = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.UpdatePartyMemberReply} returns this
 */
proto.fm.internal.UpdatePartyMemberReply.prototype.setRevision = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.InvitePartyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.InvitePartyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.InvitePartyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.InvitePartyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    inviterCharacterId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    targetCharacterName: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.InvitePartyRequest}
 */
proto.fm.internal.InvitePartyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.InvitePartyRequest;
  return proto.fm.internal.InvitePartyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.InvitePartyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.InvitePartyRequest}
 */
proto.fm.internal.InvitePartyRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setInviterCharacterId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setTargetCharacterName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.InvitePartyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.InvitePartyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.InvitePartyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.InvitePartyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getInviterCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getTargetCharacterName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.InvitePartyRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InvitePartyRequest} returns this
 */
proto.fm.internal.InvitePartyRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 inviter_character_id = 2;
 * @return {number}
 */
proto.fm.internal.InvitePartyRequest.prototype.getInviterCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InvitePartyRequest} returns this
 */
proto.fm.internal.InvitePartyRequest.prototype.setInviterCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string target_character_name = 3;
 * @return {string}
 */
proto.fm.internal.InvitePartyRequest.prototype.getTargetCharacterName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.InvitePartyRequest} returns this
 */
proto.fm.internal.InvitePartyRequest.prototype.setTargetCharacterName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.InvitePartyReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.InvitePartyReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.InvitePartyReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.InvitePartyReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    targetCharacterId: jspb.Message.getFieldWithDefault(msg, 3, 0),
    targetChannelId: jspb.Message.getFieldWithDefault(msg, 4, 0),
    partyId: jspb.Message.getFieldWithDefault(msg, 5, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.InvitePartyReply}
 */
proto.fm.internal.InvitePartyReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.InvitePartyReply;
  return proto.fm.internal.InvitePartyReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.InvitePartyReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.InvitePartyReply}
 */
proto.fm.internal.InvitePartyReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setTargetCharacterId(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setTargetChannelId(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPartyId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.InvitePartyReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.InvitePartyReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.InvitePartyReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.InvitePartyReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getTargetCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      3,
      f
    );
  }
  f = message.getTargetChannelId();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 5));
  if (f != null) {
    writer.writeUint32(
      5,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.InvitePartyReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.InvitePartyReply} returns this
 */
proto.fm.internal.InvitePartyReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.InvitePartyReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.InvitePartyReply} returns this
 */
proto.fm.internal.InvitePartyReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional uint32 target_character_id = 3;
 * @return {number}
 */
proto.fm.internal.InvitePartyReply.prototype.getTargetCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InvitePartyReply} returns this
 */
proto.fm.internal.InvitePartyReply.prototype.setTargetCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional uint32 target_channel_id = 4;
 * @return {number}
 */
proto.fm.internal.InvitePartyReply.prototype.getTargetChannelId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InvitePartyReply} returns this
 */
proto.fm.internal.InvitePartyReply.prototype.setTargetChannelId = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional uint32 party_id = 5;
 * @return {number}
 */
proto.fm.internal.InvitePartyReply.prototype.getPartyId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.InvitePartyReply} returns this
 */
proto.fm.internal.InvitePartyReply.prototype.setPartyId = function(value) {
  return jspb.Message.setField(this, 5, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.fm.internal.InvitePartyReply} returns this
 */
proto.fm.internal.InvitePartyReply.prototype.clearPartyId = function() {
  return jspb.Message.setField(this, 5, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.fm.internal.InvitePartyReply.prototype.hasPartyId = function() {
  return jspb.Message.getField(this, 5) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.DenyPartyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.DenyPartyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.DenyPartyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DenyPartyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    worldId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    deniedCharacterId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    inviterName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    action: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.DenyPartyRequest}
 */
proto.fm.internal.DenyPartyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.DenyPartyRequest;
  return proto.fm.internal.DenyPartyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.DenyPartyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.DenyPartyRequest}
 */
proto.fm.internal.DenyPartyRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setWorldId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setDeniedCharacterId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setInviterName(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setAction(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.DenyPartyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.DenyPartyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.DenyPartyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DenyPartyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorldId();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getDeniedCharacterId();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getInviterName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAction();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
};


/**
 * optional uint32 world_id = 1;
 * @return {number}
 */
proto.fm.internal.DenyPartyRequest.prototype.getWorldId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.DenyPartyRequest} returns this
 */
proto.fm.internal.DenyPartyRequest.prototype.setWorldId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint32 denied_character_id = 2;
 * @return {number}
 */
proto.fm.internal.DenyPartyRequest.prototype.getDeniedCharacterId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.DenyPartyRequest} returns this
 */
proto.fm.internal.DenyPartyRequest.prototype.setDeniedCharacterId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string inviter_name = 3;
 * @return {string}
 */
proto.fm.internal.DenyPartyRequest.prototype.getInviterName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.fm.internal.DenyPartyRequest} returns this
 */
proto.fm.internal.DenyPartyRequest.prototype.setInviterName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional uint32 action = 4;
 * @return {number}
 */
proto.fm.internal.DenyPartyRequest.prototype.getAction = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.fm.internal.DenyPartyRequest} returns this
 */
proto.fm.internal.DenyPartyRequest.prototype.setAction = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.fm.internal.DenyPartyReply.prototype.toObject = function(opt_includeInstance) {
  return proto.fm.internal.DenyPartyReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.fm.internal.DenyPartyReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DenyPartyReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    ok: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    errorCode: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.fm.internal.DenyPartyReply}
 */
proto.fm.internal.DenyPartyReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.fm.internal.DenyPartyReply;
  return proto.fm.internal.DenyPartyReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.fm.internal.DenyPartyReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.fm.internal.DenyPartyReply}
 */
proto.fm.internal.DenyPartyReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOk(value);
      break;
    case 2:
      var value = /** @type {!proto.fm.internal.PartyErrorCode} */ (reader.readEnum());
      msg.setErrorCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.fm.internal.DenyPartyReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.fm.internal.DenyPartyReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.fm.internal.DenyPartyReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.fm.internal.DenyPartyReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOk();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getErrorCode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional bool ok = 1;
 * @return {boolean}
 */
proto.fm.internal.DenyPartyReply.prototype.getOk = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.fm.internal.DenyPartyReply} returns this
 */
proto.fm.internal.DenyPartyReply.prototype.setOk = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional PartyErrorCode error_code = 2;
 * @return {!proto.fm.internal.PartyErrorCode}
 */
proto.fm.internal.DenyPartyReply.prototype.getErrorCode = function() {
  return /** @type {!proto.fm.internal.PartyErrorCode} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.fm.internal.PartyErrorCode} value
 * @return {!proto.fm.internal.DenyPartyReply} returns this
 */
proto.fm.internal.DenyPartyReply.prototype.setErrorCode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * @enum {number}
 */
proto.fm.internal.SessionErrorCode = {
  SESSION_NONE: 0,
  SESSION_UNKNOWN: 1,
  SESSION_ALREADY_LOGGED_IN: 2,
  SESSION_NOT_FOUND: 3,
  SESSION_LOGOUT_FAILED: 4
};

/**
 * @enum {number}
 */
proto.fm.internal.SessionDisconnectSource = {
  SESSION_DISCONNECT_SOURCE_UNSPECIFIED: 0,
  SESSION_DISCONNECT_SOURCE_LOGIN_SERVER: 1,
  SESSION_DISCONNECT_SOURCE_GAME_SERVER: 2
};

/**
 * @enum {number}
 */
proto.fm.internal.PartyErrorCode = {
  NONE: 0,
  UNKNOWN: 1,
  ALREADY_IN_PARTY: 2,
  CHARACTER_NOT_FOUND: 3,
  PARTY_NOT_FOUND: 4,
  PARTY_FULL: 5,
  NOT_IN_PARTY: 6,
  NOT_PARTY_LEADER: 7,
  TARGET_NOT_IN_PARTY: 8,
  TARGET_ALREADY_LEADER: 9,
  CANNOT_EXPEL_SELF: 10,
  INVITE_EXPIRED_OR_INVALID: 11,
  TARGET_OFFLINE: 12,
  INVITER_NOT_IN_PARTY: 13,
  CANNOT_INVITE_SELF: 14,
  TARGET_ALREADY_IN_PARTY: 15
};

goog.object.extend(exports, proto.fm.internal);
