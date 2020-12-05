import {Injectable, OnDestroy} from '@angular/core';
import {BehaviorSubject, interval, Observable, Subject} from 'rxjs';
import {catchError, map, pluck, startWith, switchMap, tap} from 'rxjs/operators';
import {BackendService} from './api/backend.service';
import {ResourceType, SessionResponse, UserInfoResponse} from './api/models';
import {Router} from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AuthService implements OnDestroy {

  constructor(private backend: BackendService, private router: Router) {
  }

  liveData$ = this.backend.events$.pipe(
    catchError(error => {
      console.error(error);
      return [];
    }),
    tap({
        error: error => console.log('[Live component] Error:', error),
        complete: () => console.log('[Live component] Connection Closed'),
        next: value => console.log('[Live component] Next: ', value),
      }
    )
  );

  loginSubscription = interval(600000).pipe(
    startWith([null]),
    switchMap(a => this.backend.getSession())
  ).subscribe(value => {
    this.sessionSubject.next(value);
  });

  private sessionSubject = new BehaviorSubject<SessionResponse>(undefined);
  public session$ = this.sessionSubject.asObservable();
  public authUserId$ = this.session$.pipe(pluck('id'));

  public checkLoggedIn() {
    // this.loginSubject.next(true);
  }

  public goToMyProfile() {
    if (this.sessionSubject.value !== undefined) {
      this.router.navigateByUrl('/users/' + this.sessionSubject.value.id);
    }
  }

  public goToMyResource(resourceId: string, resourceType: ResourceType) {
    if (this.sessionSubject.value !== undefined) {
      const typeUrl = resourceType === ResourceType.Offer ? 'offers' : 'needs';
      const url = '/users/' + this.sessionSubject.value.id + '/' + typeUrl + '/' + resourceId;
      this.router.navigateByUrl(url);
    }
  }

  public getUserAuthInfo(): Observable<UserInfoResponse> {
    return this.session$;
  }

  public getUserAuthId(): Observable<string> {
    return this.authUserId$;
  }

  ngOnDestroy(): void {
    if (this.loginSubscription) {
      this.loginSubscription.unsubscribe();
    }
  }

}
