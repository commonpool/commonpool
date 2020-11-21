import {Component, EventEmitter, Input, OnDestroy, Output} from '@angular/core';
import {AcceptInvitationRequest, Membership} from '../../api/models';
import {AuthService} from '../../auth.service';
import {combineLatest, Observable, ReplaySubject, Subject, Subscription} from 'rxjs';
import {map, startWith, switchMap, tap} from 'rxjs/operators';
import {BackendService} from '../../api/backend.service';

@Component({
  selector: 'app-accept-invitation',
  template: `
    <ng-container *ngIf="membership$ | async; let membership">
      <ng-container *ngIf="isMyInvitation$ | async">
        <button
          class="btn btn-primary btn-sm mr-2"
          [disabled]="pending$ | async"
          (click)="acceptInvitation(membership.userId, membership.groupId)"
        >
          Accept to join <b>{{membership.groupName}}</b>
        </button>
      </ng-container>
    </ng-container>
  `,
})
export class AcceptInvitationComponent implements OnDestroy {

  private readonly membershipSubject: Subject<Membership>;
  private readonly acceptInvitationSubject: Subject<AcceptInvitationRequest>;
  private readonly pendingSubject: Subject<boolean>;
  private readonly errorSubject: Subject<any>;

  public readonly membership$: Observable<Membership>;
  public readonly isMyInvitation$: Observable<boolean>;
  public readonly pending$: Observable<boolean>;
  public readonly error$: Observable<any>;

  public readonly acceptInvitationSubscription: Subscription;

  @Output()
  accepted: EventEmitter<Membership> = new EventEmitter<Membership>();

  constructor(private auth: AuthService, private backend: BackendService) {

    this.errorSubject = new ReplaySubject();
    this.error$ = this.errorSubject.asObservable();

    this.pendingSubject = new ReplaySubject();
    this.pending$ = this.pendingSubject.asObservable().pipe(startWith(false));

    this.membershipSubject = new ReplaySubject<Membership>(1);
    this.membership$ = this.membershipSubject.asObservable();

    this.isMyInvitation$ = combineLatest([this.membership$, this.auth.getUserAuthId()]).pipe(
      map(([membership, userAuthId]) => membership.userId === userAuthId)
    );
    this.acceptInvitationSubject = new Subject<AcceptInvitationRequest>();
    this.acceptInvitationSubscription = this.acceptInvitationSubject.asObservable().pipe(
      tap(() => this.pendingSubject.next(true)),
      tap(() => this.errorSubject.next(undefined)),
      switchMap((request) => this.backend.acceptInvitation(request)),
    ).subscribe(response => {
      this.pendingSubject.next(false);
      this.accepted.next(response.membership);
    }, err => {
      this.pendingSubject.next(false);
      this.errorSubject.next(err);
    }, () => {
      this.pendingSubject.next(false);
    });
  }

  @Input()
  set membership(value: Membership) {
    this.membershipSubject.next(value);
  }

  acceptInvitation(userId: string, groupId: string) {
    this.acceptInvitationSubject.next(new AcceptInvitationRequest(userId, groupId));
  }

  ngOnDestroy(): void {
    if (this.acceptInvitationSubscription) {
      this.acceptInvitationSubscription.unsubscribe();
    }
  }

}
