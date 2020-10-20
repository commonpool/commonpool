import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {pluck, switchMap} from 'rxjs/operators';
import {BackendService} from '../api/backend.service';
import {BehaviorSubject, Subject} from 'rxjs';
import {SearchResourceRequest, UserInfoResponse} from '../api/models';

@Component({
  selector: 'app-user-profile',
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {

  userId$ = this.route.params.pipe(pluck('id'));
  userResources$ = this.userId$.pipe(
    switchMap(id => this.backend.searchResources(new SearchResourceRequest(undefined, undefined, id, 10, 0))),
    pluck('resources')
  );
  userInfo$ = this.userId$.pipe(
    switchMap(id => this.backend.getUserInfo(id))
  );

  constructor(private route: ActivatedRoute, private backend: BackendService) {
  }

  ngOnInit(): void {
  }

}
