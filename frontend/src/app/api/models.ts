import {HttpResponse} from '@angular/common/http';
import {formatDistanceToNow} from 'date-fns';

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
    if (res.body.code && res.body.message && res.body.statusCode) {
      return new ErrorResponse(res.body.message, res.body.code, res.body.statusCode);
    }
    return new ErrorResponse(res.body, '', res.status);
  }
}

export class UserInfoResponse {
  constructor(public id: string, public username: string) {
  }
}


export class GetThreadsRequest {
  constructor(public skip: number, public take: number) {
  }
}

export class Thread {
  constructor(
    public hasUnreadMessages: boolean,
    public topicId: string,
    public lastChars: string,
    public lastMessageAt: string,
    public lastMessageUserId: string,
    public lastMessageUsername: string) {
  }

  static from(thread: Thread) {
    return new Thread(
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
    this.sentAtDistance = formatDistanceToNow(Date.parse(sentAt), {addSuffix: true});
  }

  public sentAtDistance: string;

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


