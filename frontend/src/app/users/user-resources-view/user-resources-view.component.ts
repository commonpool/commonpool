import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {map, pluck, switchMap} from 'rxjs/operators';
import {ExtendedResource, ResourceType, SearchResourceRequest} from '../../api/models';
import {combineLatest, Observable} from 'rxjs';
import {UserInfoService} from '../user-info.service';

@Component({
  selector: 'app-user-resources-view',
  templateUrl: './user-resources-view.component.html',
  styleUrls: ['./user-resources-view.component.css']
})
export class UserResourcesViewComponent {

  resourceType$: Observable<ResourceType>;
  isOffers$: Observable<boolean>;
  isNeeds$: Observable<boolean>;
  userId$: Observable<string>;
  resources$: Observable<ExtendedResource[]>;
  isMyProfile$: Observable<boolean>;

  constructor(private route: ActivatedRoute, private backend: BackendService, private userService: UserInfoService) {
    this.isMyProfile$ = this.userService.getIsMyProfile();
    this.resourceType$ = this.route.url.pipe(
      map(u => u[0].path === 'offers' ? ResourceType.Offer : ResourceType.Request)
    );
    this.isOffers$ = this.resourceType$.pipe(map(r => r === ResourceType.Offer));
    this.isNeeds$ = this.resourceType$.pipe(map(r => r === ResourceType.Request));
    this.userId$ = this.route.parent.params.pipe(pluck('id'));
    this.resources$ = combineLatest([
      this.userId$,
      this.resourceType$
    ]).pipe(
      switchMap(([userId, resourceType]) => {
        return this.backend.searchResources(new SearchResourceRequest(undefined, resourceType, userId, 10, 0));
      }),
      pluck('resources')
    );
  }

}
