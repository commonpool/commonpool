import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {pluck, switchMap} from 'rxjs/operators';
import {GetMyMembershipsRequest} from '../../api/models';
import {AuthService} from '../../auth.service';

@Component({
  selector: 'app-user-groups-view',
  templateUrl: './user-groups-view.component.html',
  styleUrls: ['./user-groups-view.component.css']
})
export class UserGroupsViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService, private auth: AuthService) {
  }

  userId$ = this.route.parent.params.pipe(pluck('id'));
  groups$ = this.userId$.pipe(
    switchMap(id => this.backend.getMyMemberships(new GetMyMembershipsRequest()))
  );
  authUser$ = this.auth.session$.pipe(pluck('id'))

  ngOnInit(): void {
  }

}
