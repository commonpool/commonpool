import {HttpResponse} from '@angular/common/http';

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
    public exchangeValue: number,
    public timeSensitivity: number,
    public necessityLevel: number) {
  }
}

export class SearchResourcesResponse {
  constructor(public resources: Resource[], public totalCount: number, public  take: number, public skip: number) {

  }
}

export class SearchResourceRequest {
  constructor(public query: string, public type: ResourceType, public take: number, public skip: number) {
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
