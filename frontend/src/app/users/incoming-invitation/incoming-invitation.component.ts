import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {AcceptInvitationRequest, DeclineInvitationRequest, Membership} from '../../api/models';
import {BackendService} from '../../api/backend.service';
import {pluck} from 'rxjs/operators';
import {Observable} from 'rxjs';

@Component({
  selector: 'app-incoming-invitation',
  templateUrl: './incoming-invitation.component.html',
  styleUrls: ['./incoming-invitation.component.css']
})
export class IncomingInvitationComponent implements OnInit {

  constructor(private backend: BackendService) {
  }

  @Input()
  isProfileOwner: boolean;

  @Input()
  membership: Membership;

  @Output()
  accepted: EventEmitter<void> = new EventEmitter<void>();

  @Output()
  declined: EventEmitter<void> = new EventEmitter<void>();

  pending = false;

  error: any = undefined;

  accept() {
    const acceptAtion = this.backend
      .acceptInvitation(new AcceptInvitationRequest(this.membership.userId, this.membership.groupId))
      .pipe(pluck('membership'));
    const afterAccept = () => this.accepted.next();
    this.acceptOrDecline(acceptAtion, afterAccept);
  }


  decline() {
    const declineAction = this.backend
      .declineInvitation(new DeclineInvitationRequest(this.membership.userId, this.membership.groupId))
      .pipe(pluck('membership'));
    const afterDecline = () => this.declined.next();
    this.acceptOrDecline(declineAction, afterDecline);
  }

  private acceptOrDecline<A>(backendAction: Observable<A>, after: () => void) {
    this.pending = true;
    this.error = undefined;
    backendAction.subscribe(res => {
      this.pending = false;
      after();
    }, err => {
      this.error = err;
      this.pending = false;
    }, () => {
      this.pending = false;
    });
  }

  ngOnInit(): void {
  }

}
