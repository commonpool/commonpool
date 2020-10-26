import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {pluck, switchMap} from 'rxjs/operators';
import {BackendService} from '../../api/backend.service';
import {GetGroupMembershipsRequest} from '../../api/models';

@Component({
  selector: 'app-group-members-view',
  templateUrl: './group-members-view.component.html',
  styleUrls: ['./group-members-view.component.css']
})
export class GroupMembersViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService) {
  }

  params$ = this.route.parent.params;
  groupId$ = this.params$.pipe(pluck('id'));
  members$ = this.groupId$.pipe(
    switchMap(id => this.backend.getGroupMemberships(new GetGroupMembershipsRequest(id)))
  );

  ngOnInit(): void {
    this.params$.subscribe(p => console.log(p))
  }

}
