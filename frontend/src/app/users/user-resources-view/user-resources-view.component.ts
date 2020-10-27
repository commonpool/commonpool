import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {pluck, switchMap} from 'rxjs/operators';
import {GetMyMembershipsRequest, SearchResourceRequest} from '../../api/models';

@Component({
  selector: 'app-user-resources-view',
  templateUrl: './user-resources-view.component.html',
  styleUrls: ['./user-resources-view.component.css']
})
export class UserResourcesViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService) {
  }

  userId$ = this.route.parent.params.pipe(pluck('id'));
  resources$ = this.userId$.pipe(
    switchMap(id => this.backend.searchResources(new SearchResourceRequest(undefined, undefined, id, 10, 0)))
  );

  ngOnInit(): void {
  }

}
