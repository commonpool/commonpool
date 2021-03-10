import {HttpResponse} from '@angular/common/http';
import {formatDistanceToNow, format} from 'date-fns';

export enum CallType {
  Offer = 'offer',
  Request = 'request'
}

export enum ResourceType {
  Object = 'object',
  Service = 'service'
}

export enum OfferItemType {
  CreditTransfer = 'transfer_credits',
  ProvideService = 'provide_service',
  BorrowResource = 'borrow_resource',
  ResourceTransfer = 'transfer_resource'
}

export class SharedWithOutput {
  public constructor(public groupId: string, public groupName: string) {
  }

  public static from(sharedWithOutput: SharedWithOutput): SharedWithOutput {
    return new SharedWithOutput(sharedWithOutput.groupId, sharedWithOutput.groupName);
  }
}

export class SharedWithInput {
  public constructor(public groupId: string) {
  }

  public static from(sharedWithInput: SharedWithInput): SharedWithInput {
    return new SharedWithInput(sharedWithInput.groupId);
  }
}

export class Resource {
  constructor(
    public resourceId: string,
    public createdBy: string,
    public createdByVersion: number,
    public createdByName: string,
    public createdAt: string,
    public updatedBy: string,
    public updatedByVersion: number,
    public updatedByName: string,
    public updatedAt: string,
    public groupSharingCount: number,
    public version: number,
    public owner: Target,
    public info: ResourceInfo,
    public sharings: SharedWithOutput[]
  ) {
  }

  public static from(r: Resource): Resource {
    return new Resource(
      r.resourceId,
      r.createdBy,
      r.createdByVersion,
      r.createdByName,
      r.createdAt,
      r.updatedBy,
      r.updatedByVersion,
      r.updatedByName,
      r.updatedAt,
      r.groupSharingCount,
      r.version,
      Target.from(r.owner),
      ResourceInfo.from(r.info),
      r.sharings.map(s => SharedWithOutput.from(s)));
  }
}

export class ExtendedResource extends Resource {
  public createdAtDistance: string;
  public createdAtDistanceAgo: string;

  constructor(r: Resource) {
    super(
      r.resourceId,
      r.createdBy,
      r.createdByVersion,
      r.createdByName,
      r.createdAt,
      r.updatedBy,
      r.updatedByVersion,
      r.updatedByName,
      r.updatedAt,
      r.groupSharingCount,
      r.version,
      r.owner,
      r.info,
      r.sharings);
    this.createdAtDistance = formatDistanceToNow(Date.parse(r.createdAt));
    this.createdAtDistanceAgo = formatDistanceToNow(Date.parse(r.createdAt), {addSuffix: true});
  }
}

export class ResourceValue {
  public valueType: string;

  public constructor(
    public timeValueFrom: number,
    public timeValueTo: number
  ) {
    this.valueType = 'from_to_duration';
  }

  public static from(r: ResourceValue): ResourceValue {
    return new ResourceValue(r.timeValueFrom, r.timeValueTo);
  }
}

export class ResourceInfoUpdate {
  public constructor(
    public name: string,
    public description: string,
    public value: ResourceValue,
  ) {

  }

  public static from(r: ResourceInfoUpdate): ResourceInfoUpdate {
    return new ResourceInfoUpdate(r.name, r.description, ResourceValue.from(r.value));
  }
}

export class ResourceInfo extends ResourceInfoUpdate {
  public constructor(
    public name: string,
    public description: string,
    public value: ResourceValue,
    public callType: CallType,
    public resourceType: ResourceType,
  ) {
    super(name, description, value);
  }

  public static from(r: ResourceInfo): ResourceInfo {
    return new ResourceInfo(r.name, r.description, ResourceValue.from(r.value), r.callType, r.resourceType);
  }
}

export class GetResourceResponse {
  public resource: ExtendedResource;

  constructor(resource: Resource) {
    this.resource = new ExtendedResource(resource);
  }

  public static from(res: GetResourceResponse): GetResourceResponse {
    return new GetResourceResponse(Resource.from(res.resource));
  }
}

export class SearchResourcesResponse {
  public resources: ExtendedResource[];

  constructor(resources: Resource[], public totalCount: number, public  take: number, public skip: number) {
    this.resources = resources.map(r => new ExtendedResource(r));
  }
}

export class SearchResourceRequest {
  constructor(
    public query: string,
    public callType: CallType,
    public resourceType: ResourceType,
    public createdBy: string,
    public groupId,
    public take: number,
    public skip: number) {
  }
}

export class CreateResourcePayload {
  constructor(
    public info: ResourceInfo,
    public sharedWith: SharedWithInput[]) {
  }

  public static from(p: CreateResourcePayload): CreateResourcePayload {
    return new CreateResourcePayload(
      ResourceInfo.from(p.info),
      p.sharedWith ? p.sharedWith.map(w => SharedWithInput.from(w)) : []
    );
  }
}

export class UpdateResourcePayload {
  constructor(
    public info: ResourceInfoUpdate,
    public sharedWith: SharedWithInput[]) {
  }

  public static from(r: UpdateResourcePayload): UpdateResourcePayload {
    return new UpdateResourcePayload(
      ResourceInfoUpdate.from(r.info),
      r.sharedWith ? r.sharedWith.map(w => SharedWithInput.from(w)) : []
    );
  }
}

export class CreateResourceRequest {
  constructor(
    public resource: CreateResourcePayload) {
  }

  public static from(r: CreateResourceRequest): CreateResourceRequest {
    return new CreateResourceRequest(CreateResourcePayload.from(r.resource));
  }
}

export class UpdateResourceRequest {
  constructor(
    public id: string,
    public resource: UpdateResourcePayload) {
  }

  public static from(r: UpdateResourceRequest): UpdateResourceRequest {
    return new UpdateResourceRequest(r.id, UpdateResourcePayload.from(r.resource));
  }
}

export class CreateResourceResponse {
  public resource: ExtendedResource;

  constructor(resource: Resource) {
    this.resource = new ExtendedResource(resource);
  }
}

export class UpdateResourceResponse {
  public resource: ExtendedResource;

  constructor(resource: Resource) {
    this.resource = new ExtendedResource(resource);
  }
}

export class SessionResponse {
  constructor(public username: string, public id: string, public isAuthenticated: boolean) {
  }
}

export class ErrorResponse {

  constructor(public message: string, public code: string, statusCode: number) {
  }

  static fromHttpResponse(res: HttpResponse<any>): ErrorResponse {
    if (res?.body?.code && res?.body?.message && res?.body?.statusCode) {
      return new ErrorResponse(res.body.message, res.body.code, res.body.statusCode);
    }
    if (res.body) {
      return new ErrorResponse(res.body, '', res.status);
    }
    return new ErrorResponse(res.statusText, '', res.status);

  }
}

export class UserInfoResponse {
  constructor(public id: string, public username: string) {
  }

  static from(res: UserInfoResponse): UserInfoResponse {
    return new UserInfoResponse(res.id, res.username);
  }
}

export class UsersInfoResponse {
  constructor(public users: UserInfoResponse[]) {
  }

  static from(res: UsersInfoResponse): UsersInfoResponse {
    return new UsersInfoResponse(res.users.map(u => UserInfoResponse.from(u)));
  }
}

export class SearchUsersQuery {
  constructor(public query: string, public take: number, public skip: number) {
  }
}

export class GetSubscriptionsRequest {
  constructor(public skip: number, public take: number) {
  }
}

export enum ChannelType {
  Group = 0,
  Conversation = 1
}

export class Subscription {
  constructor(
    public channelId: string,
    public userId: string,
    public createdAt: string,
    public updatedAt: string,
    public lastMessageAt: string,
    public lastTimeRead: string,
    public lastMessageChars: string,
    public lastMessageUserId: string,
    public lastMessageUsername: string,
    public name: string,
    public type: ChannelType) {
  }

  static from(s: Subscription) {
    return new Subscription(
      s.channelId,
      s.userId,
      s.createdAt,
      s.updatedAt,
      s.lastMessageAt,
      s.lastTimeRead,
      s.lastMessageChars,
      s.lastMessageUserId,
      s.lastMessageUsername,
      s.name,
      s.type);
  }
}

export class GetChannelMembershipsResponse {
  constructor(public subscriptions: Subscription[]) {
  }

  static from(response: GetChannelMembershipsResponse) {
    return new GetChannelMembershipsResponse(response.subscriptions.map(t => Subscription.from(t)));
  }
}

export class GetMessagesRequest {
  constructor(public skip: number, public take: number, public topic: string) {
  }
}

export enum MessageType {
  NormalMessage = 'message'
}

export enum MessageSubType {
  UserMessage = 'user',
  BotMessage = 'bot'
}

export class Message {
  constructor(
    public id: string,
    public channelId: string,
    public messageType: MessageType,
    public messageSubType: MessageSubType,
    public sentById: string,
    public sentByUsername: string,
    public sentAt: string,
    public text: string,
    public blocks: Block[],
    public attachments: Attachment[],
    public visibleToUser: string) {
    const date = new Date(Date.parse(this.sentAt));
    this.sentAtDistance = formatDistanceToNow(date, {addSuffix: true});
    this.sentAtHour = format(date, 'hh');
    this.sentAtDate = date;
  }

  public sentAtDistance: string;
  public sentAtHour: string;
  public sentAtDate: Date;

  static from(m: Message) {
    return new Message(
      m.id,
      m.channelId,
      m.messageType,
      m.messageSubType,
      m.sentById,
      m.sentByUsername,
      m.sentAt,
      m.text,
      m.blocks ? m.blocks.map(b => Block.from(b)) : undefined,
      m.attachments ? m.attachments.map(a => Attachment.from(a)) : undefined,
      m.visibleToUser
    );
  }
}

export class GetMessagesResponse {
  constructor(public messages: Message[]) {
  }

  static from(response: GetMessagesResponse): GetMessagesResponse {
    return new GetMessagesResponse(response.messages.map(m => Message.from(m)));
  }
}

export class SendMessageRequest {
  constructor(public topic: string, public content: string) {
  }
}

export class InquireAboutResourceRequest {
  constructor(public resourceId: string, public content: string) {
  }
}

export class SendOfferRequest {
  constructor(public offer: SendOfferRequestPayload) {
  }

  public static from(req: SendOfferRequest): SendOfferRequest {
    return new SendOfferRequest(SendOfferRequestPayload.from(req.offer));
  }
}

export class GetOffersRequest {
}

export class SendOfferRequestPayload {
  constructor(public items: SendOfferRequestItem[], public message: string, public groupId: string) {
  }

  public static from(req: SendOfferRequestPayload): SendOfferRequestPayload {
    return new SendOfferRequestPayload(req.items.map(i => SendOfferRequestItem.from(i)), req.message, req.groupId);
  }
}

export enum TargetType {
  User = 'user',
  Group = 'group'
}

export class Target {
  constructor(public type: TargetType, public userId: string | undefined, public  groupId: string | undefined) {
  }

  public static from(target: Target): Target {
    if (!target) {
      return undefined;
    }
    return new Target(target.type, target.userId, target.groupId);
  }
}

export class SendOfferRequestItem {
  constructor(
    public to: Target,
    public type: OfferItemType,
    public from: Target | undefined,
    public resourceId: string | undefined,
    public duration: string | undefined,
    public amount: string | undefined) {
  }

  static from(req: SendOfferRequestItem) {
    return new SendOfferRequestItem(Target.from(req.to), req.type, Target.from(req.from), req.resourceId, req.duration, req.amount);
  }
}

export class GetOfferRequest {
  constructor(id: string) {
  }
}

export class GetOfferResponse {
  constructor(public offer: Offer) {
  }

  public static from(res: GetOfferResponse): GetOfferResponse {
    return new GetOfferResponse(Offer.from(res.offer));
  }
}

export class GetOffersResponse {
  constructor(public offers: Offer[]) {
  }

  public static from(res: GetOffersResponse): GetOffersResponse {
    return new GetOffersResponse(res.offers.map(o => Offer.from(o)));
  }
}

export class SendOfferResponse {
  constructor(public offer: Offer) {
  }

  public static from(res: SendOfferResponse): SendOfferResponse {
    return new SendOfferResponse(Offer.from(res.offer));
  }
}

export class AcceptOfferRequest {
  constructor(public id: string) {
  }
}

export class AcceptOfferResponse {
  constructor(public offer: Offer) {
  }

  public static from(o: AcceptOfferResponse): AcceptOfferResponse {
    return new AcceptOfferResponse(Offer.from(o.offer));
  }
}

export class DeclineOfferRequest {
  constructor(public id: string) {
  }
}

export class DeclineOfferReponse {
  constructor(public offer: Offer) {
  }

  public static from(o: DeclineOfferReponse): DeclineOfferReponse {
    return new DeclineOfferReponse(Offer.from(o.offer));
  }
}

export enum OfferStatus {
  PendingOffer = 0,
  AcceptedOffer = 1,
  CanceledOffer = 2,
  DeclinedOffer = 3,
  ExpiredOffer = 4,
  CompletedOffer = 5,
}

export class OfferItem {
  constructor(
    public id: string,
    public to: Target,
    public type: OfferItemType,
    public from: Target | undefined,
    public resourceId: string | undefined,
    public duration: number | undefined,
    public amount: number | undefined,
    public receivingApprovers: string[],
    public givingApprovers: string[],
    public giverApproved: boolean,
    public receiverApproved: boolean,
    public serviceGivenConfirmation: boolean,
    public serviceReceivedConfirmation: boolean,
    public itemTaken: boolean,
    public itemGiven: boolean,
    public itemReturnedBack: boolean,
    public itemReceivedBack: boolean) {
  }

  public static from(res: OfferItem): OfferItem {
    return new OfferItem(
      res.id,
      Target.from(res.to),
      res.type,
      Target.from(res.from),
      res.resourceId,
      res.duration,
      res.amount,
      res.receivingApprovers,
      res.givingApprovers,
      res.giverApproved,
      res.receiverApproved,
      res.serviceGivenConfirmation,
      res.serviceReceivedConfirmation,
      res.itemTaken,
      res.itemGiven,
      res.itemReturnedBack,
      res.itemReceivedBack);
  }
}

export class Offer {
  constructor(
    public id: string,
    public createdAt: string,
    public completedAt: string,
    public status: OfferStatus,
    public authorId: string,
    public authorUsername: string,
    public items: OfferItem[]) {
  }

  public static from(o: Offer): Offer {
    return new Offer(
      o.id,
      o.createdAt,
      o.completedAt,
      o.status,
      o.authorId,
      o.authorUsername,
      o.items.map(i => OfferItem.from(i)));
  }
}

export class Group {
  constructor(
    public id: string,
    public createdAt: string,
    public name: string,
    public description: string
  ) {
  }

  public static from(g: Group): Group {
    return new Group(
      g.id,
      g.createdAt,
      g.name,
      g.description);
  }
}

export class Membership {
  constructor(
    public version: number,
    public groupId: string,
    public groupName: string,
    public userId: string,
    public isOwner: boolean,
    public isAdmin: boolean,
    public isMember: boolean,
    public groupConfirmed: boolean,
    public groupConfirmedBy: boolean,
    public groupConfirmedAt: string,
    public userConfirmed: boolean,
    public userConfirmedAt: string,
    public status: MembershipStatus,
    public userVersion: number,
    public userName: string,
    public createdBy: string,
    public createdByName: string,
    public createdByVersion: number,
    public createdAt: string
  ) {
    this.createdAtDate = new Date(Date.parse(createdAt));
    this.createdAtDistance = formatDistanceToNow(this.createdAtDate, {addSuffix: true});
  }

  createdAtDate: Date;
  createdAtDistance: string;

  public static from(m: Membership): Membership {
    return new Membership(
      m.version,
      m.groupId,
      m.groupName,
      m.userId,
      m.isOwner,
      m.isAdmin,
      m.isMember,
      m.groupConfirmed,
      m.groupConfirmedBy,
      m.groupConfirmedAt,
      m.userConfirmed,
      m.userConfirmedAt,
      m.status,
      m.userVersion,
      m.userName,
      m.createdBy,
      m.createdByName,
      m.createdByVersion,
      m.createdAt
    );
  }
}

export class CreateGroupRequest {
  constructor(public name: string, description: string) {
  }
}

export class CreateGroupResponse {
  constructor(public group: Group) {
  }

  public static from(r: CreateGroupResponse): CreateGroupResponse {
    return new CreateGroupResponse(Group.from(r.group));
  }
}

export class GetGroupRequest {
  constructor(public id: string) {
  }
}

export class GetGroupResponse {
  constructor(public group: Group) {
  }

  public static from(g: GetGroupResponse): GetGroupResponse {
    return new GetGroupResponse(Group.from(g.group));
  }
}

export class InviteUserRequest {
  constructor(public userId: string, public groupId: string) {
  }
}

export class InviteUserResponse {
  constructor(public membership: Membership) {
  }

  public static from(i: InviteUserResponse): InviteUserResponse {
    return new InviteUserResponse(Membership.from(i.membership));
  }
}

export class ExcludeUserRequest {
  constructor(public userId: string, public groupId: string) {
  }
}

export class ExcludeUserResponse {
  constructor(public membership: Membership) {
  }

  public static from(i: ExcludeUserResponse): ExcludeUserResponse {
    return new ExcludeUserResponse(Membership.from(i.membership));
  }
}

export enum PermissionType {
  MemberPermission,
  AdminPermission
}

export class GrantPermissionRequest {
  constructor(public userId: string, public groupId: string, public permission: PermissionType) {
  }
}

export class GrantPermissionResponse {
  constructor(public membership: Membership) {
  }

  public static from(i: GrantPermissionResponse): GrantPermissionResponse {
    return new GrantPermissionResponse(Membership.from(i.membership));
  }
}

export class RevokePermissionRequest {
  constructor(public userId: string, public groupId: string, public permission: PermissionType) {
  }
}

export class RevokePermissionResponse {
  constructor(public membership: Membership) {
  }

  public static from(i: RevokePermissionResponse): RevokePermissionResponse {
    return new RevokePermissionResponse(Membership.from(i.membership));
  }
}

export class GetMyMembershipsRequest {

}

export class GetMyMembershipsResponse {
  constructor(public memberships: Membership[]) {
  }

  public static from(i: GetMyMembershipsResponse): GetMyMembershipsResponse {
    return new GetMyMembershipsResponse(i.memberships.map(m => Membership.from(m)));
  }
}

export enum MembershipStatus {
  PendingGroupMembershipStatus,
  PendingUserMembershipStatus,
  ApprovedMembershipStatus,

}

export class GetUserMembershipsRequest {
  constructor(public userId: string, public membershipStatus?: MembershipStatus) {
  }
}

export class GetUserMembershipsResponse {
  constructor(public memberships: Membership[]) {
  }

  public static from(i: GetUserMembershipsResponse): GetUserMembershipsResponse {
    return new GetUserMembershipsResponse(i.memberships.map(m => Membership.from(m)));
  }
}

export class GetGroupMembershipsRequest {
  constructor(public groupId: string, public membershipStatus?: MembershipStatus) {
  }
}

export class GetGroupMembershipsResponse {
  constructor(public memberships: Membership[]) {
  }

  public static from(r: GetGroupMembershipsResponse): GetGroupMembershipsResponse {
    return new GetGroupMembershipsResponse(r.memberships.map(m => Membership.from(m)));
  }
}

export class GetMembershipRequest {
  constructor(public userId: string, public groupId: string) {
  }
}

export class GetMembershipResponse {
  constructor(public membership: Membership) {
  }

  public static from(r: GetMembershipResponse): GetMembershipResponse {
    return new GetMembershipResponse(Membership.from(r.membership));
  }
}

export class GetUsersForGroupInvitePickerRequest {
  constructor(public skip: number, public take: number, public query: string, public groupId: string) {
  }
}

export class GetUsersForGroupInvitePickerResponse {
  constructor(public users: UserInfoResponse[], public skip: number, public take: number) {
  }

  static from(res: GetUsersForGroupInvitePickerResponse): GetUsersForGroupInvitePickerResponse {
    return new GetUsersForGroupInvitePickerResponse(res.users.map(u => UserInfoResponse.from(u)), res.skip, res.take);
  }
}

export class AcceptInvitationRequest {
  constructor(public userId: string, public groupId: string) {
  }
}

export class AcceptInvitationResponse {
  constructor(public membership: Membership) {
  }

  public static from(r: AcceptInvitationResponse): AcceptInvitationResponse {
    return new AcceptInvitationResponse(Membership.from(r.membership));
  }
}

export class DeclineInvitationRequest {
  constructor(public userId: string, public groupId: string) {
  }
}

export class DeclineInvitationResponse {
  constructor(public membership: Membership) {
  }

  public static from(r: DeclineInvitationResponse): DeclineInvitationResponse {
    return new DeclineInvitationResponse(Membership.from(r.membership));
  }
}

export class LeaveGroupRequest {
  constructor(public userId: string, public groupId: string) {
  }
}

export class LeaveGroupResponse {
  constructor(public membership: Membership) {
  }

  public static from(r: LeaveGroupResponse): LeaveGroupResponse {
    return new LeaveGroupResponse(Membership.from(r.membership));
  }
}

export enum SurfaceType {
  Modals = 'modals',
  Messages = 'messages',
  Home = 'home',
  Tabs = 'tabs',
}

export enum BlockType {
  Actions = 'actions',
  Context = 'context',
  Divider = 'divider',
  File = 'file',
  Header = 'header',
  Image = 'image',
  Input = 'input',
  Section = 'section'
}

export enum ElementType {
  ButtonElement = 'button',
  PlainTextInputElement = 'plain_text_input',
  ImageElement = 'image',
  CheckboxesElement = 'checkboxes',
  DatepickerElement = 'datepicker',
  RadioButtonsElement = 'radio_buttons',
}

export enum ButtonStyle {
  Primary = 'primary',
  Danger = 'danger'
}

export enum TextType {
  PlainTextType = 'plain_text',
  MarkdownTextType = 'mrkdwn'
}

export enum FileSource {
  Remote = 'remote'
}

export const BlockTypeCompatibility: { [key: string]: SurfaceType[] } = {};
BlockTypeCompatibility[BlockType.Actions] = [SurfaceType.Modals, SurfaceType.Messages, SurfaceType.Home, SurfaceType.Tabs];
BlockTypeCompatibility[BlockType.Context] = [SurfaceType.Modals, SurfaceType.Messages, SurfaceType.Home, SurfaceType.Tabs];
BlockTypeCompatibility[BlockType.Divider] = [SurfaceType.Modals, SurfaceType.Messages, SurfaceType.Home, SurfaceType.Tabs];
BlockTypeCompatibility[BlockType.File] = [SurfaceType.Messages];
BlockTypeCompatibility[BlockType.Header] = [SurfaceType.Modals, SurfaceType.Messages, SurfaceType.Home, SurfaceType.Tabs];
BlockTypeCompatibility[BlockType.Image] = [SurfaceType.Modals, SurfaceType.Messages, SurfaceType.Home, SurfaceType.Tabs];
BlockTypeCompatibility[BlockType.Input] = [SurfaceType.Modals, SurfaceType.Home, SurfaceType.Tabs];
BlockTypeCompatibility[BlockType.Section] = [SurfaceType.Modals, SurfaceType.Messages, SurfaceType.Home, SurfaceType.Tabs];

export const BlockElementCompatibility: { [key: string]: BlockType[] } = {};
BlockElementCompatibility[ElementType.ButtonElement] = [BlockType.Section, BlockType.Actions];
BlockElementCompatibility[ElementType.CheckboxesElement] = [BlockType.Section, BlockType.Actions, BlockType.Input];
BlockElementCompatibility[ElementType.DatepickerElement] = [BlockType.Section, BlockType.Actions, BlockType.Input];
BlockElementCompatibility[ElementType.ImageElement] = [BlockType.Section, BlockType.Context];
BlockElementCompatibility[ElementType.PlainTextInputElement] = [BlockType.Section, BlockType.Actions, BlockType.Input];
BlockElementCompatibility[ElementType.RadioButtonsElement] = [BlockType.Section, BlockType.Actions, BlockType.Input];

export class TextObject {
  constructor(public type: TextType, public value: string, public emoji?: boolean) {
  }

  public static from(textObject: TextObject): TextObject {
    return new TextObject(textObject.type, textObject.value, textObject.emoji);
  }
}

export class OptionObject {
  public constructor(public text: TextObject, public value: string, public description: TextObject) {
  }

  public static from(optionObject: OptionObject): OptionObject {
    return new OptionObject(TextObject.from(optionObject.text), optionObject.value, TextObject.from(optionObject.description));
  }
}

export class OptionGroupObject {
  public constructor(public text: TextObject, public options: OptionObject[]) {
  }

  public static from(optionGroupObject: OptionGroupObject): OptionGroupObject {
    return new OptionGroupObject(TextObject.from(optionGroupObject.text), optionGroupObject.options);
  }
}

export class Block {
  public constructor(
    public type: BlockType,
    public text?: TextObject,
    public elements?: (BlockElement | TextObject)[],
    public imageUrl?: string,
    public altText?: string,
    public title?: TextObject,
    public fields?: TextObject[],
    public accessory?: BlockElement,
    public blockId?: string,
    public externalId?: string,
    public source?: FileSource) {
  }

  public static from(b: Block): Block {
    return new Block(
      b.type,
      b.text ? TextObject.from(b.text) : undefined,
      b.elements ? b.elements.map(e => {
        if (e.type === TextType.MarkdownTextType || e.type === TextType.PlainTextType) {
          return TextObject.from(e as TextObject);
        } else {
          return BlockElement.from(e as BlockElement);
        }
      }) : undefined,
      b.imageUrl,
      b.altText,
      b.title ? TextObject.from(b.title) : undefined,
      b.fields ? b.fields.map(f => TextObject.from(f)) : undefined,
      b.accessory ? BlockElement.from(b.accessory) : undefined,
      b.blockId,
      b.externalId,
      b.source
    );
  }
}

export class ActionsBlock {
  public type: BlockType = BlockType.Actions;

  public constructor(public elements: (BlockElement | TextObject)[], public blockId?: string) {
  }

  public static from(b: Block): ActionsBlock {
    if (b.type !== BlockType.Actions) {
      throw new Error('invalid block type');
    }
    return new ActionsBlock(b.elements, b.blockId);
  }
}

export class ContextBlock {
  public type: BlockType = BlockType.Context;

  public constructor(public elements: (BlockElement | TextObject)[], public blockId?: string) {
  }

  public static from(b: Block): ContextBlock {
    if (b.type !== BlockType.Context) {
      throw new Error('invalid block type');
    }
    return new ContextBlock(b.elements, b.blockId);
  }
}

export class DividerBlock {
  public type: BlockType = BlockType.Divider;

  public constructor(public blockId?: string) {
  }

  public static from(b: Block): DividerBlock {
    if (b.type !== BlockType.Divider) {
      throw new Error('invalid block type');
    }
    return new DividerBlock(b.blockId);
  }
}

export class FileBlock {
  public type: BlockType = BlockType.File;

  public constructor(public externalId: string, public source: FileSource, public blockId?: string) {
  }

  public static from(b: Block): DividerBlock {
    if (b.type !== BlockType.File) {
      throw new Error('invalid block type');
    }
    return new FileBlock(b.externalId, b.source, b.blockId);
  }
}

export class HeaderBlock {
  public type: BlockType = BlockType.Header;

  public constructor(public text: TextObject, public blockId?: string) {
  }

  public static from(b: Block): HeaderBlock {
    if (b.type !== BlockType.Header) {
      throw new Error('invalid block type');
    }
    return new HeaderBlock(b.text, b.blockId);
  }
}

export class ImageBlock {
  public type: BlockType = BlockType.Image;

  public constructor(public imageUrl: string, public altText: string, public title?: TextObject, public blockId?: string) {
  }

  public static from(b: Block): ImageBlock {
    if (b.type !== BlockType.Image) {
      throw new Error('invalid block type');
    }
    return new ImageBlock(b.imageUrl, b.altText, b.title ? TextObject.from(b.title) : undefined, b.blockId);
  }
}

export class SectionBlock {
  public type: BlockType = BlockType.Section;

  public constructor(
    public text: TextObject,
    public fields?: TextObject[],
    public accessory?: BlockElement,
    public blockId?: string) {
  }

  public static from(b: Block): SectionBlock {
    if (b.type !== BlockType.Section) {
      throw new Error('invalid block type');
    }
    return new SectionBlock(b.text, b.fields, b.accessory, b.blockId);
  }
}

export class BlockElement {
  public constructor(
    public type: ElementType,
    public text?: TextObject,
    public actionId?: string,
    public url?: string,
    public value?: string,
    public style?: ButtonStyle,
    public confirm?: boolean,
    public options?: OptionObject[],
    public initialOptions?: OptionObject[],
    public placeholder?: TextObject,
    public initialDate?: string,
    public imageUrl?: string,
    public altText?: string
  ) {
  }

  public static from(b: BlockElement): BlockElement {
    return new BlockElement(
      b.type,
      b.text ? TextObject.from(b.text) : undefined,
      b.actionId,
      b.url,
      b.value,
      b.style,
      b.confirm,
      b.options?.map(o => OptionObject.from(o)),
      b.initialOptions?.map(o => OptionObject.from(o)),
      b.placeholder ? TextObject.from(b.placeholder) : undefined,
      b.initialDate,
      b.imageUrl,
      b.altText
    );
  }
}

export class ButtonElement {
  public type: ElementType = ElementType.ButtonElement;

  public constructor(
    public text: TextObject,
    public style?: ButtonStyle,
    public actionId?: string,
    public url?: string,
    public value?: string,
    public confirm?: boolean) {
  }

  public static from(b: BlockElement): ButtonElement {
    if (b.type !== ElementType.ButtonElement) {
      throw new Error('invalid block type');
    }
    return new ButtonElement(TextObject.from(b.text), b.style, b.actionId, b.url, b.value, b.confirm);
  }
}

export class CheckboxesElement {
  public type: ElementType = ElementType.CheckboxesElement;

  public constructor(
    text: TextObject,
    public actionId: string,
    public options: OptionObject[],
    public initialOptions?: OptionObject[],
    public confirm?: boolean) {
  }

  public static from(b: BlockElement): CheckboxesElement {
    if (b.type !== ElementType.CheckboxesElement) {
      throw new Error('invalid block type');
    }
    return new CheckboxesElement(
      b.text,
      b.actionId,
      b.options?.map(o => OptionObject.from(o)),
      b.initialOptions?.map(o => OptionObject.from(o)),
      b.confirm);
  }
}

export class DatePickerElement {
  public type: ElementType = ElementType.DatepickerElement;

  public constructor(public actionId: string, public placeholder?: TextObject, public initialDate?: string, public confirm?: boolean) {
  }

  public static from(b: BlockElement): DatePickerElement {
    if (b.type !== ElementType.DatepickerElement) {
      throw new Error('invalid block type');
    }
    return new DatePickerElement(
      b.actionId,
      b.placeholder ? TextObject.from(b.placeholder) : undefined,
      b.initialDate,
      b.confirm
    );
  }
}

export class ImageElement {
  public type: ElementType = ElementType.ImageElement;

  public constructor(public imageUrl: string, public altText: string) {
  }

  public static from(b: BlockElement): ImageElement {
    if (b.type !== ElementType.ImageElement) {
      throw new Error('invalid block type');
    }
    return new ImageElement(b.imageUrl, b.altText);
  }
}

export class RadioButtonsElement {
  public type: ElementType = ElementType.RadioButtonsElement;

  public constructor(
    public actionId: string,
    public options: OptionObject[],
    public initialOptions?: OptionObject[],
    public confirm?: boolean) {
  }

  public static from(b: BlockElement): RadioButtonsElement {
    if (b.type !== ElementType.RadioButtonsElement) {
      throw new Error('invalid block type');
    }
    return new RadioButtonsElement(
      b.actionId,
      b.options?.map(o => OptionObject.from(o)),
      b.initialOptions?.map(o => OptionObject.from(o)),
      b.confirm);
  }
}

export class Attachment {
  public constructor(public color: string, public blocks: Block[]) {
  }

  public static from(a: Attachment): Attachment {
    return new Attachment(a.color, a.blocks.map(b => Block.from(b)));
  }
}

export class ElementState {
  public constructor(
    public type: ElementType,
    public selectedDate?: string,
    public selectedTime?: string,
    public value?: string,
    public selectedOption?: OptionObject,
    public selectedOptions?: OptionObject[]
  ) {
  }

  public static from(e: ElementState): ElementState {
    return new ElementState(
      e.type,
      e.selectedDate,
      e.selectedTime,
      e.value,
      e.selectedOption ? OptionObject.from(e.selectedOption) : undefined,
      e.selectedOptions ? e.selectedOptions.map(o => OptionObject.from(o)) : undefined
    );
  }
}

export class SubmitAction extends ElementState {
  constructor(
    public blockId: string,
    public actionId: string,
    public type: ElementType,
    public selectedDate?: string,
    public selectedTime?: string,
    public value?: string,
    public selectedOption?: OptionObject,
    public selectedOptions?: OptionObject[]) {
    super(type, selectedDate, selectedTime, value, selectedOption, selectedOptions);
  }

  public static from(e: SubmitAction): SubmitAction {
    return new SubmitAction(
      e.blockId,
      e.actionId,
      e.type,
      e.selectedDate,
      e.selectedTime,
      e.value,
      e.selectedOption ? OptionObject.from(e.selectedOption) : undefined,
      e.selectedOptions ? e.selectedOptions.map(o => OptionObject.from(o)) : undefined);
  }
}

export type SubmitActionState = {
  [blockId: string]: { [actionId: string]: ElementState }
};

export class SubmitInteractionPayload {
  public constructor(public messageId: string, public actions: SubmitAction[], public state: SubmitActionState) {
  }

  public static from(p: SubmitInteractionPayload): SubmitInteractionPayload {
    const newState: SubmitActionState = {};
    for (const blockKey in p.state) {
      if (p.state.hasOwnProperty(blockKey)) {
        for (const actionKey in p.state[blockKey]) {
          if (p.state[blockKey].hasOwnProperty(actionKey)) {
            if (!newState[blockKey]) {
              newState[blockKey] = {};
            }
            if (!newState[blockKey][actionKey]) {
              newState[blockKey][actionKey] = {} as ElementState;
            }
            newState[blockKey][actionKey] = ElementState.from(p.state[blockKey][actionKey]);
          }
        }
      }
    }
    return new SubmitInteractionPayload(
      p.messageId,
      p.actions ? p.actions.map(a => SubmitAction.from(a)) : undefined,
      newState);
  }
}

export class SubmitInteractionRequest {
  public constructor(public payload: SubmitInteractionPayload) {
  }
}

export enum EventType {
  MessageEvent = 'message'
}

export enum EventSubType {
  MessageChanged = 'message_changed',
  MessageDeleted = 'message_deleted'
}

export class Event {
  public constructor(
    public type: EventType,
    public subType: EventSubType,
    public channel: string,
    public user: string,
    public id: string,
    public timestamp: string,
    public text: string) {
  }

  public static from(e: Event): Event {
    return new Event(
      e.type,
      e.subType,
      e.channel,
      e.user,
      e.id,
      e.timestamp,
      e.text);
  }
}

export class GetTradingHistoryRequest {
  public constructor(public userIds: string[]) {
  }

  public static from(r: GetTradingHistoryRequest): GetTradingHistoryRequest {
    return new GetTradingHistoryRequest(r.userIds);
  }
}

export class TradingHistoryEntry {
  public constructor(
    public timestamp: string,
    public fromUserId: string,
    public toUserId: string,
    public fromUsername: string,
    public toUsername: string,
    public resourceId: string | undefined,
    public timeAmountSeconds: number | undefined) {
  }

  public static from(t: TradingHistoryEntry): TradingHistoryEntry {
    return new TradingHistoryEntry(
      t.timestamp,
      t.fromUserId,
      t.toUserId,
      t.fromUsername,
      t.toUsername,
      t.resourceId,
      t.timeAmountSeconds
    );
  }
}

export class GetTradingHistoryResponse {
  public constructor(public entries: TradingHistoryEntry[]) {

  }

  public static from(r: GetTradingHistoryResponse): GetTradingHistoryResponse {
    return new GetTradingHistoryResponse(
      r.entries ? r.entries.map(e => TradingHistoryEntry.from(e)) : []
    );
  }
}

export class OfferGroupOrUserPickerItem {
  public constructor(public type: TargetType, public name: string, public userId: string | undefined, public groupId: string | undefined) {

  }

  public static from(i: OfferGroupOrUserPickerItem): OfferGroupOrUserPickerItem {
    return new OfferGroupOrUserPickerItem(i.type, i.name, i.userId, i.groupId);
  }

}

export class OfferGroupOrUserPickerResult {
  public constructor(public items: OfferGroupOrUserPickerItem[]) {

  }

  public static from(o: OfferGroupOrUserPickerResult): OfferGroupOrUserPickerResult {
    return new OfferGroupOrUserPickerResult(o.items ? o.items.map(i => OfferGroupOrUserPickerItem.from(i)) : undefined);
  }

}

export class OfferItemTargetRequest {
  public constructor(
    public type: OfferItemType,
    public groupId: string,
    public fromType: TargetType | undefined,
    public fromId: string | undefined,
    public toType: TargetType | undefined,
    public toId: string | undefined) {
  }
}

export class ConfirmServiceProvidedRequest {
  public constructor(public offerItemId: string) {
  }
}

export class ConfirmResourceTransferred {
  public constructor(public offerItemId: string) {
  }
}

export class ConfirmResourceBorrowed {
  public constructor(public offerItemId: string) {
  }
}

export class ConfirmBorrowedResourceReturned {
  public constructor(public offerItemId: string) {
  }
}


