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
  GetChannelMembershipsResponse,
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
  GetGroupResponse,
  GetGroupMembershipsRequest,
  GetGroupMembershipsResponse,
  GetUsersForGroupInvitePickerRequest,
  GetUsersForGroupInvitePickerResponse,
  GetUserMembershipsRequest,
  GetUserMembershipsResponse,
  GetMembershipRequest,
  GetMembershipResponse,
  AcceptInvitationRequest,
  AcceptInvitationResponse,
  DeclineInvitationRequest,
  DeclineInvitationResponse,
  LeaveGroupRequest, LeaveGroupResponse, SubmitInteractionRequest, Event
} from './models';

import {Observable, of, Subject, throwError} from 'rxjs';
import {
  HttpClient,
  HttpEvent,
  HttpHandler,
  HttpHeaders,
  HttpInterceptor,
  HttpRequest,
  HttpResponse
} from '@angular/common/http';
import {catchError, map, retry, switchAll, tap} from 'rxjs/operators';
import {environment} from '../../environments/environment';
import {WebSocketSubject} from 'rxjs/internal-compatibility';
import {webSocket} from 'rxjs/webSocket';

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
        const bla = evt as any;
        if (bla?.body?.meta?.redirectTo) {
          setTimeout(() => {
            window.location = bla.body.meta.redirectTo;
          }, 1000);
        }
      }),
      catchError((err: any) => {
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

  private socket$: WebSocketSubject<any>;
  private messagesSubject$ = new Subject<Event>();
  public events$ = this.messagesSubject$;

  public connect(): void {

    console.log('connecting');

    const socketObserver = new Observable(observer => {
      try {

        const subject = webSocket(`${environment.apiUrl
          .replace('https://', 'wss://')
          .replace('http://', 'wss://')
        }/api/v1/ws`);

        //const subject = webSocket(`ws://localhost:8585/api/v1/ws`);

        const subscription = subject.asObservable()
          .subscribe(data =>
              observer.next(data),
            error => observer.error(error),
            () => observer.complete());
        return () => {
          if (!subscription.closed) {
            subscription.unsubscribe();
          }
        };
      } catch (error) {
        observer.error(error);
      }
    });

    const messages = socketObserver.pipe(
      retry(Infinity),
      map(a => Event.from(a as any)))
      .subscribe((e) => {
        this.messagesSubject$.next(e);
      });

  }

  sendMessageWs(msg: any) {
    this.socket$.next(msg);
  }

  close() {
    this.socket$.complete();
  }

  login() {
    this.http.post(`${environment.apiUrl}/api/v1/auth/login`, undefined, {
      observe: 'response'
    }).subscribe(() => {
      console.log('logging in');
    });
  }

  logout() {
    this.http.post(`${environment.apiUrl}/api/v1/auth/logout`, undefined, {
      observe: 'response'
    }).subscribe(() => {
      console.log('logging out');
    });
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
    if (request.groupId !== undefined) {
      params.group_id = request.groupId;
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
        return GetResourceResponse.from(body);
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

  getSubscriptions(skip: number, take: number): Observable<GetChannelMembershipsResponse> {
    const params: { [key: string]: string } = {};
    if (skip) {
      params.skip = skip.toString();
    }
    if (take) {
      params.take = take.toString();
    }
    return this.http.get(`${environment.apiUrl}/api/v1/chat/subscriptions`, {
      observe: 'response',
      params
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetChannelMembershipsResponse.from(res.body as GetChannelMembershipsResponse);
      })
    );
  }

  getMessages(channelId: string, before: number, take: number): Observable<GetMessagesResponse> {
    const params: { [key: string]: string } = {};
    if (before) {
      params.before = before.toString();
    }
    if (take) {
      params.take = take.toString();
    }
    if (channelId) {
      params.channel = channelId;
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

  getUserMemberships(request: GetUserMembershipsRequest): Observable<GetUserMembershipsResponse> {
    const params: { [key: string]: string } = {};
    if (request.membershipStatus !== undefined) {
      params.status = request.membershipStatus.toString();
    }
    return this.http.get(`${environment.apiUrl}/api/v1/users/${request.userId}/memberships`, {
      observe: 'response',
      params
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetUserMembershipsResponse.from(res.body as GetUserMembershipsResponse);
      })
    );
  }

  getGroupMemberships(request: GetGroupMembershipsRequest): Observable<GetGroupMembershipsResponse> {
    const params: { [key: string]: string } = {};
    if (request.membershipStatus !== undefined) {
      params.status = request.membershipStatus.toString();
    }
    return this.http.get(`${environment.apiUrl}/api/v1/groups/${request.groupId}/memberships`, {
      observe: 'response',
      params
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetGroupMembershipsResponse.from(res.body as GetGroupMembershipsResponse);
      })
    );
  }

  getMembership(request: GetMembershipRequest): Observable<GetMembershipResponse> {
    return this.http.get(`${environment.apiUrl}/api/v1/groups/${request.groupId}/memberships/${request.userId}`, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetMembershipResponse.from(res.body as GetMembershipResponse);
      })
    );
  }

  getUsersForGroupInvitePicker(query: GetUsersForGroupInvitePickerRequest): Observable<GetUsersForGroupInvitePickerResponse> {
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
    return this.http.get(`${environment.apiUrl}/api/v1/groups/${query.groupId}/invite-member-picker`, {
      observe: 'response',
      params
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return GetUsersForGroupInvitePickerResponse.from(res.body as GetUsersForGroupInvitePickerResponse);
      })
    );
  }

  inviteUser(request: InviteUserRequest): Observable<InviteUserResponse> {
    return this.http.post(`${environment.apiUrl}/api/v1/memberships`, {
      userId: request.userId,
      groupId: request.groupId
    }, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 201) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return InviteUserResponse.from(res.body as InviteUserResponse);
      })
    );
  }

  acceptInvitation(request: AcceptInvitationRequest): Observable<AcceptInvitationResponse> {
    return this.inviteUser(new InviteUserRequest(request.userId, request.groupId));
  }

  declineInvitation(request: DeclineInvitationRequest): Observable<DeclineInvitationResponse> {
    return this.http.request('DELETE', `${environment.apiUrl}/api/v1/memberships`, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      }),
      body: {
        userId: request.userId,
        groupId: request.groupId
      },
    }).pipe(
      map((res: HttpResponse<object>) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return DeclineInvitationResponse.from(res.body as DeclineInvitationResponse);
      })
    );
  }

  leaveGroup(request: LeaveGroupRequest): Observable<LeaveGroupResponse> {
    return this.declineInvitation(new DeclineInvitationRequest(request.userId, request.groupId));
  }

  submitMessageInteraction(request: SubmitInteractionRequest): Observable<any> {
    return this.http.post(`${environment.apiUrl}/api/v1/chat/interaction`, request, {
      observe: 'response'
    }).pipe(
      map((res) => {
        if (res.status !== 200) {
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        return res;
      })
    );
  }

}
