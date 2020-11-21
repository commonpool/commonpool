import {Injectable, OnDestroy} from '@angular/core';
import {BehaviorSubject, Observable, Subject} from 'rxjs';
import {pluck, switchMap} from 'rxjs/operators';
import {BackendService} from './api/backend.service';
import {SessionResponse, UserInfoResponse} from './api/models';
import {Router} from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AuthService implements OnDestroy {

  constructor(private backend: BackendService, private router: Router) {
  }

  loginSubject = new Subject();
  loginSubscription = this.loginSubject.asObservable()
    .pipe(
      switchMap(a => this.backend.getSession())
    ).subscribe(value => {
      this.sessionSubject.next(value);
    });

  private sessionSubject = new BehaviorSubject<SessionResponse>(undefined);
  public session$ = this.sessionSubject.asObservable();
  public authUserId$ = this.session$.pipe(pluck('id'));

  public checkLoggedIn() {
    this.loginSubject.next(true);
  }

  public goToMyProfile() {
    if (this.sessionSubject.value !== undefined) {
      this.router.navigateByUrl('/users/' + this.sessionSubject.value.id);
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
