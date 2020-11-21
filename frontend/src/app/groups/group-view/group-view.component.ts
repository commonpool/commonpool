import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {pluck, switchMap} from 'rxjs/operators';
import {BackendService} from '../../api/backend.service';
import {GetGroupRequest} from '../../api/models';
import {GroupService} from '../group.service';

@Component({
  selector: 'app-group-view',
  templateUrl: './group-view.component.html',
  styleUrls: ['./group-view.component.css']
})
export class GroupViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService, private groupService: GroupService) {
    this.groupId$.subscribe(id => this.groupService.setGroupId(id));
  }

  groupId$ = this.route.params.pipe(pluck('id'));
  group$ = this.groupId$.pipe(
    switchMap(id => this.backend.getGroup(new GetGroupRequest(id))),
    pluck('group')
  );

  ngOnInit(): void {
  }

}
