import {Injectable, OnInit} from '@angular/core';
import {combineLatest, Observable, of, ReplaySubject, Subject} from 'rxjs';
import {BackendService} from '../api/backend.service';
import {map, shareReplay, startWith, switchMap, tap} from 'rxjs/operators';
import {UserInfoResponse} from '../api/models';
import {AuthService} from '../auth.service';

@Injectable({
  providedIn: 'root'
})
export class UserInfoService {

  private readonly userIdSubject: Subject<string>;
  private readonly userId$: Observable<string>;
  private readonly userInfo$: Observable<UserInfoResponse>;
  private readonly isMyProfile$: Observable<boolean>;

  constructor(private backend: BackendService, private auth: AuthService) {
    this.userIdSubject = new ReplaySubject<string>(1);
    this.userId$ = this.userIdSubject.asObservable().pipe();
    this.userInfo$ = this.userId$.pipe(
      switchMap((userId) => {
        if (userId) {
          return this.backend.getUserInfo(userId);
        } else {
          return of(undefined as UserInfoResponse);
        }
      }),
      shareReplay()
    );
    this.isMyProfile$ = combineLatest([this.getUserId(), this.auth.authUserId$]).pipe(
      map(([userId, authUserId]) => userId === authUserId),
      shareReplay(),
      startWith(false)
    );
  }

  setUserId(userId: string) {
    this.userIdSubject.next(userId);
  }

  public getUserId(): Observable<string> {
    return this.userId$;
  }

  public getIsMyProfile(): Observable<boolean> {
    return this.isMyProfile$;
  }

  public getUserInfo(): Observable<UserInfoResponse> {
    return this.userInfo$;
  }
}
