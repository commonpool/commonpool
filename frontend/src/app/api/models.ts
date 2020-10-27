import {HttpResponse} from '@angular/common/http';
import {formatDistanceToNow, format} from 'date-fns';

export enum ResourceType {
  Offer = 0,
  Request = 1
}

export class Resource {
  constructor(
    public id: string,
    public summary: string,
    public description: string,
    public type: ResourceType,
    public valueInHoursFrom: number,
    public valueInHoursTo: number,
    public createdBy: string,
    public createdById: string,
    public createdAt: string
  ) {
  }

  public static from(res: Resource): Resource {
    return new Resource(
      res.id,
      res.summary,
      res.description,
      res.type,
      res.valueInHoursFrom,
      res.valueInHoursTo,
      res.createdBy,
      res.createdById,
      res.createdAt);
  }
}

export class ExtendedResource extends Resource {
  public createdAtDistance: string;
  public createdAtDistanceAgo: string;

  constructor(resource: Resource) {
    super(
      resource.id,
      resource.summary,
      resource.description,
      resource.type,
      resource.valueInHoursFrom,
      resource.valueInHoursTo,
      resource.createdBy,
      resource.createdById,
      resource.createdAt);
    this.createdAtDistance = formatDistanceToNow(Date.parse(resource.createdAt));
    this.createdAtDistanceAgo = formatDistanceToNow(Date.parse(resource.createdAt), {addSuffix: true});
  }
}

export class GetResourceResponse {
  public resource: Resource;

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
  constructor(public query: string, public type: ResourceType, public createdBy: string, public take: number, public skip: number) {
  }
}

export class CreateResourcePayload {
  constructor(
    public summary: string,
    public description: string,
    public resourceType: ResourceType,
    public valueInHoursFrom: number,
    public valueInHoursTo: number) {
  }
}

export class UpdateResourcePayload {
  constructor(
    public summary: string,
    public description: string,
    public resourceType: ResourceType,
    public valueInHoursFrom: number,
    public valueInHoursTo: number) {
  }
}

export class CreateResourceRequest {
  constructor(
    public resource: CreateResourcePayload) {
  }
}

export class UpdateResourceRequest {
  constructor(
    public id: string,
    public resource: UpdateResourcePayload) {
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

export class GetThreadsRequest {
  constructor(public skip: number, public take: number) {
  }
}

export class Thread {
  constructor(
    public title: string,
    public hasUnreadMessages: boolean,
    public topicId: string,
    public lastChars: string,
    public lastMessageAt: string,
    public lastMessageUserId: string,
    public lastMessageUsername: string) {
  }

  static from(thread: Thread) {
    return new Thread(
      thread.title,
      thread.hasUnreadMessages,
      thread.topicId,
      thread.lastChars,
      thread.lastMessageAt,
      thread.lastMessageUserId,
      thread.lastMessageUsername);
  }
}

export class GetThreadsResponse {
  constructor(public threads: Thread[]) {
  }

  static from(response: GetThreadsResponse) {
    return new GetThreadsResponse(response.threads.map(t => Thread.from(t)));
  }
}

export class GetMessagesRequest {
  constructor(public skip: number, public take: number, public topic: string) {
  }
}

export class Message {
  constructor(
    public content: string,
    public id: string,
    public sentAt: string,
    public sentBy: string,
    public sentByUsername: string,
    public topicId: string,
    public sentByMe: boolean) {
    const date = new Date(Date.parse(this.sentAt));
    this.sentAtDistance = formatDistanceToNow(date, {addSuffix: true});
    this.sentAtHour = format(date, 'hh');
    this.sentAtDate = date;
  }

  public sentAtDistance: string;
  public sentAtHour: string;
  public sentAtDate: Date;

  static from(message: Message) {
    return new Message(
      message.content,
      message.id,
      message.sentAt,
      message.sentBy,
      message.sentByUsername,
      message.topicId,
      message.sentByMe);
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
  constructor(public items: SendOfferRequestItem[]) {
  }

  public static from(req: SendOfferRequestPayload): SendOfferRequestPayload {
    return new SendOfferRequestPayload(req.items.map(i => SendOfferRequestItem.from(i)));
  }
}

export enum OfferItemType {
  ResourceItem = 0,
  TimeItem = 1
}

export class SendOfferRequestItem {
  constructor(public from: string, public to: string, public type: OfferItemType, public resourceId: string, public timeInSeconds: number) {
  }

  static from(req: SendOfferRequestItem) {
    return new SendOfferRequestItem(req.from, req.to, req.type, req.resourceId, req.timeInSeconds);
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
  ExpiredOffer = 4
}

export class OfferItem {
  constructor(
    public id: string,
    public fromUserId: string,
    public toUserId: string,
    public type: OfferItemType,
    public resourceId: string,
    public timeInSeconds: number) {
  }

  public static from(res: OfferItem): OfferItem {
    return new OfferItem(res.id, res.fromUserId, res.toUserId, res.type, res.resourceId, res.timeInSeconds);
  }
}

export enum Decision {
  PendingDecision = 0,
  AcceptedDecision = 1,
  DeclinedDecision = 2
}

export class OfferDecision {
  constructor(public offerId: string, public userId: string, public decision: Decision) {
  }

  public static from(res: OfferDecision): OfferDecision {
    return new OfferDecision(res.offerId, res.userId, res.decision);
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
    public items: OfferItem[],
    public decisions: OfferDecision[]) {
  }

  public static from(o: Offer): Offer {
    return new Offer(
      o.id,
      o.createdAt,
      o.completedAt,
      o.status,
      o.authorId,
      o.authorUsername,
      o.items.map(i => OfferItem.from(i)),
      o.decisions.map(d => OfferDecision.from(d)));
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
    public userId: string,
    public groupId: string,
    public isAdmin: boolean,
    public isMember: boolean,
    public isOwner: boolean,
    public groupConfirmed: boolean,
    public userConfirmed: boolean,
    public createdAt: string,
    public isDeactivated: boolean,
    public groupName: string,
    public userName: string
  ) {
    this.createdAtDate = new Date(Date.parse(createdAt));
    this.createdAtDistance = formatDistanceToNow(this.createdAtDate, {addSuffix: true});
  }

  createdAtDate: Date;
  createdAtDistance: string;

  public static from(m: Membership): Membership {
    return new Membership(
      m.userId,
      m.groupId,
      m.isAdmin,
      m.isMember,
      m.isOwner,
      m.groupConfirmed,
      m.userConfirmed,
      m.createdAt,
      m.isDeactivated,
      m.groupName,
      m.userName
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
  ApprovedMembershipStatus,
  PendingStatus,
  PendingGroupMembershipStatus,
  PendingUserMembershipStatus
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
  constructor(public id: string) {
  }
}

export class GetGroupMembershipsResponse {
  constructor(public memberships: Membership[]) {
  }

  public static from(r: GetGroupMembershipsResponse): GetGroupMembershipsResponse {
    return new GetGroupMembershipsResponse(r.memberships.map(m => Membership.from(m)));
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
