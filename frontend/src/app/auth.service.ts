import {Injectable} from '@angular/core';
import {BehaviorSubject, Subject} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {BackendService} from './api/backend.service';
import {SessionResponse} from './api/models';
import {Router} from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private backend: BackendService, private router: Router) {
  }

  loginRequestSubject$ = new Subject();
  reqSub = this.loginRequestSubject$.asObservable()
    .pipe(
      switchMap(a => this.backend.getSession())
    ).subscribe(value => {
      this.sessionSubject.next(value);
    });

  private sessionSubject = new BehaviorSubject<SessionResponse>(undefined);
  public session$ = this.sessionSubject.asObservable();

  public checkLoggedIn() {
    this.loginRequestSubject$.next(true);
  }

  public goToMyProfile() {
    if (this.sessionSubject.value !== undefined) {
      this.router.navigateByUrl('/users/' + this.sessionSubject.value.id);
    }
  }


}
