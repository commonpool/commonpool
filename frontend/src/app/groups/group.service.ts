import {Injectable} from '@angular/core';
import {combineLatest, Observable, ReplaySubject, Subject} from 'rxjs';
import {GetMembershipRequest, Group, Membership} from '../api/models';
import {BackendService} from '../api/backend.service';
import {pluck, shareReplay, switchMap} from 'rxjs/operators';
import {AuthService} from '../auth.service';

@Injectable({
  providedIn: 'root'
})
export class GroupService {

  private permissionMap: { [userKey: string]: { [groupId: string]: Observable<Membership> } } = {};

  private readonly groupIdSubject: Subject<string>;
  private readonly groupId$: Observable<string>;
  private readonly groupSubject: Subject<Group>;
  private readonly group$: Observable<Group>;
  private readonly myMembership$: Observable<Membership>;

  constructor(private backend: BackendService, private auth: AuthService) {
    this.groupIdSubject = new ReplaySubject(1);
    this.groupId$ = this.groupIdSubject.asObservable();
    this.groupSubject = new ReplaySubject(1);
    this.group$ = this.groupSubject.asObservable();
    this.myMembership$ = combineLatest([
      this.groupId$,
      this.auth.getUserAuthId()
    ]).pipe(
      switchMap(([groupId, userAuthId]) => this.getPermission(userAuthId, groupId))
    );
  }

  public setGroupId(groupId: string) {
    this.groupIdSubject.next(groupId);
  }

  public getMyMembership(): Observable<Membership> {
    return this.myMembership$;
  }

  public getPermission(userId: string, groupId: string): Observable<Membership> {
    this.refreshPermission(userId, groupId);
    return this.permissionMap[userId][groupId];
  }

  public refreshPermission(userId: string, groupId: string) {
    if (!this.permissionMap[userId]) {
      this.permissionMap[userId] = {};
    }
    if (!this.permissionMap[userId][groupId]) {
      this.permissionMap[userId][groupId] = this.backend
        .getMembership(new GetMembershipRequest(userId, groupId))
        .pipe(
          pluck('membership'),
          shareReplay()
        );
    }
  }

}
