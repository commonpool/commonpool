import {Injectable} from '@angular/core';
import {
  CreateResourceRequest,
  CreateResourceResponse,
  ErrorResponse,
  SessionResponse,
  SearchResourceRequest,
  SearchResourcesResponse,
  GetResourceResponse,
  UserInfoResponse,
  UpdateResourceResponse,
  UpdateResourceRequest,
  GetThreadsResponse,
  GetMessagesResponse,
  UsersInfoResponse,
  SearchUsersQuery,
  SendOfferRequest,
  SendOfferResponse,
  GetOfferRequest,
  GetOfferResponse,
  GetOffersResponse,
  AcceptOfferRequest,
  AcceptOfferResponse,
  DeclineOfferRequest,
  DeclineOfferReponse,
  CreateGroupRequest,
  CreateGroupResponse,
  InviteUserRequest,
  InviteUserResponse,
  GetMyMembershipsRequest,
  GetMyMembershipsResponse,
  GetGroupRequest,
  GetGroupResponse
} from './models';

import {Observable, of, throwError} from 'rxjs';
import {HttpClient, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from '@angular/common/http';
import {catchError, map, tap} from 'rxjs/operators';
import {environment} from '../../environments/environment';

@Injectable()
export class AppHttpInterceptor implements HttpInterceptor {
  constructor() {
  }

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    return next.handle(req).pipe(
      tap(evt => {

      }),
      catchError((err: any) => {
        console.log(err);
        if (err?.error?.meta?.redirectTo) {
          setTimeout(() => {
            window.location = err.error.meta.redirectTo;
          }, 1000);
        }
        return of(err);
      }));
  }
}

@Injectable({
  providedIn: 'root'
})
export class BackendService {

  constructor(private http: HttpClient) {

  }

  createResource(request: CreateResourceRequest): Observable<CreateResourceResponse> {
    return this.http.post<CreateResourceResponse>(`${environment.apiUrl}/api/v1/resources`, request, {
      observe: 'response',
    }).pipe(
      map((res) => {
        if (res.status > 399 || res.body === null) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        const body = res.body as CreateResourceResponse;
        return new CreateResourceResponse(body.resource);
      })
    );
  }

  updateResource(request: UpdateResourceRequest): Observable<CreateResourceResponse> {
    return this.http.put<UpdateResourceResponse>(`${environment.apiUrl}/api/v1/resources/` + request.id, request, {
      observe: 'response',
    }).pipe(
      map((res) => {
        if (res.status > 399 || res.body === null) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        const body = res.body as UpdateResourceResponse;
        return new UpdateResourceResponse(body.resource);
      })
    );
  }

  searchResources(request: SearchResourceRequest): Observable<SearchResourcesResponse> {
    const params: any = {};
    if (request.createdBy !== undefined) {
      params.created_by = request.createdBy;
    }
    if (request.take !== undefined) {
      params.take = request.take.toString();
    }
    if (request.skip !== undefined) {
      params.skip = request.skip.toString();
    }
    if (request.type !== undefined) {
      params.type = request.type.toString();
    }
    if (request.query !== undefined) {
      params.query = request.query;
    }
    return this.http.get<SearchResourcesResponse>(`${environment.apiUrl}/api/v1/resources`, {
      params,
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status > 399 || res.body === null) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        const body = res.body as SearchResourcesResponse;
        return new SearchResourcesResponse(body.resources, body.totalCount, body.take, body.skip);
      })
    );
  }

  getResource(id: string): Observable<GetResourceResponse> {
    return this.http.get(`${environment.apiUrl}/api/v1/resources/${id}`, {
      observe: 'response',
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        const body = res.body as GetResourceResponse;
        return new GetResourceResponse(body.resource);
      })
    );
  }

  getSession(): Observable<SessionResponse> {
    return this.http.get(`${environment.apiUrl}/api/v1/meta/who-am-i`, {
      observe: 'response',
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        const body = res.body as SessionResponse;
        return new SessionResponse(body.username, body.id, body.isAuthenticated);
      })
    );
  }

  getUserInfo(id: string): Observable<UserInfoResponse> {
    return this.http.get(`${environment.apiUrl}/api/v1/users/` + id, {
      observe: 'response',
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return UserInfoResponse.from(res.body as UserInfoResponse);
      })
    );
  }

  searchUsers(query: SearchUsersQuery): Observable<UsersInfoResponse> {
    const params: { [key: string]: string } = {};
    if (query.skip) {
      params.skip = query.skip.toString();
    }
    if (query.take) {
      params.take = query.take.toString();
    }
    if (query.query) {
      params.query = query.query.toString();
    }
    return this.http.get(`${environment.apiUrl}/api/v1/users`, {
      observe: 'response',
      params
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return UsersInfoResponse.from(res.body as UsersInfoResponse);
      })
    );
  }

  getThreads(skip: number, take: number): Observable<GetThreadsResponse> {
    const params: { [key: string]: string } = {};
    if (skip) {
      params.skip = skip.toString();
    }
    if (take) {
      params.take = take.toString();
    }
    return this.http.get(`${environment.apiUrl}/api/v1/chat/threads`, {
      observe: 'response',
      params
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetThreadsResponse.from(res.body as GetThreadsResponse);
      })
    );
  }

  getMessages(topic: string, skip: number, take: number): Observable<GetMessagesResponse> {
    const params: { [key: string]: string } = {};
    if (skip) {
      params.skip = skip.toString();
    }
    if (take) {
      params.take = take.toString();
    }
    if (topic) {
      params.topic = topic;
    }
    return this.http.get(`${environment.apiUrl}/api/v1/chat/messages`, {
      observe: 'response',
      params
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetMessagesResponse.from(res.body as GetMessagesResponse);
      })
    );
  }

  inquireAboutResource(resource: string, content: string): Observable<void> {
    return this.http.post(`${environment.apiUrl}/api/v1/resources/${resource}/inquire`, {message: content}, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return;
      })
    );
  }

  sendMessage(topic: string, content: string): Observable<void> {
    return this.http.post(`${environment.apiUrl}/api/v1/chat/${topic}`, {message: content}, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 202) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return;
      })
    );
  }

  sendOffer(offer: SendOfferRequest): Observable<SendOfferResponse> {
    return this.http.post(`${environment.apiUrl}/api/v1/offers`, offer, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 202) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return SendOfferResponse.from(res.body as SendOfferResponse);
      })
    );
  }

  getOffers(offer: GetOfferRequest): Observable<GetOffersResponse> {
    return this.http.get(`${environment.apiUrl}/api/v1/offers`, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetOffersResponse.from(res.body as GetOffersResponse);
      })
    );
  }


  acceptOffer(offer: AcceptOfferRequest): Observable<AcceptOfferResponse> {
    return this.http.post(`${environment.apiUrl}/api/v1/offers/${offer.id}/accept`, undefined, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return AcceptOfferResponse.from(res.body as AcceptOfferResponse);
      })
    );
  }


  declineOffer(offer: DeclineOfferRequest): Observable<DeclineOfferReponse> {
    return this.http.post(`${environment.apiUrl}/api/v1/offers/${offer.id}/decline`, undefined, {

      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return DeclineOfferReponse.from(res.body as DeclineOfferReponse);
      })
    );
  }

  createGroup(request: CreateGroupRequest): Observable<CreateGroupResponse> {
    return this.http.post(`${environment.apiUrl}/api/v1/groups/`, request, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 201) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return CreateGroupResponse.from(res.body as CreateGroupResponse);
      })
    );
  }

  getGroup(request: GetGroupRequest): Observable<GetGroupResponse> {
    return this.http.get(`${environment.apiUrl}/api/v1/groups/${request.id}`, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetGroupResponse.from(res.body as GetGroupResponse);
      })
    );
  }

  getMyMemberships(request: GetMyMembershipsRequest): Observable<GetMyMembershipsResponse> {
    return this.http.get(`${environment.apiUrl}/api/v1/my/memberships`, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetMyMembershipsResponse.from(res.body as GetMyMembershipsResponse);
      })
    );
  }
}
