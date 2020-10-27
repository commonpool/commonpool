import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {pluck, switchMap} from 'rxjs/operators';
import {GetMyMembershipsRequest} from '../../api/models';

@Component({
  selector: 'app-user-invitations-view',
  templateUrl: './user-invitations-view.component.html',
  styleUrls: ['./user-invitations-view.component.css']
})
export class UserInvitationsViewComponent implements OnInit {


  constructor(private route: ActivatedRoute, private backend: BackendService) {
  }

  userId$ = this.route.parent.params.pipe(pluck('id'));
  groups$ = this.userId$.pipe(
    switchMap(id => this.backend.getMyMemberships(new GetMyMembershipsRequest()))
  );


  ngOnInit(): void {
  }

}
