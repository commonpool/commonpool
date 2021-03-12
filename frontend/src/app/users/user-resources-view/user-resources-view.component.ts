import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {map, pluck, switchMap, tap} from 'rxjs/operators';
import {ExtendedResource, CallType, SearchResourceRequest} from '../../api/models';
import {combineLatest, Observable} from 'rxjs';
import {UserInfoService} from '../user-info.service';

export enum AccountType {
  Group = 'group',
  User = 'user'
}

export enum ResourceTypeStr {
  Needs = 'needs',
  Offers = 'offers'
}

@Component({
  selector: 'app-user-resources-view',
  templateUrl: './user-resources-view.component.html',
  styleUrls: ['./user-resources-view.component.css']
})
export class UserResourcesViewComponent implements OnInit {

  resourceType$: Observable<CallType>;
  isOffers$: Observable<boolean>;
  isNeeds$: Observable<boolean>;
  accountId$: Observable<string>;
  resources$: Observable<ExtendedResource[]>;
  isMyProfile$: Observable<boolean>;
  ownerType$: Observable<string>;

  constructor(
    private route: ActivatedRoute,
    private backend: BackendService,
    private userService: UserInfoService
  ) {
  }

  ngOnInit(): void {

    this.ownerType$ = this.route.data.pipe(map(d => d.accountType === 'group' ? 'group' : 'user'));

    this.isMyProfile$ = combineLatest([this.userService.getIsMyProfile(), this.route.data]).pipe(
      map(([isMyProfile, data]) => data.accountType === AccountType.User && isMyProfile),
    );

    this.resourceType$ = this.route.data.pipe(
      map(d => d.resourceType === ResourceTypeStr.Offers ? CallType.Offer : d.resourceType === ResourceTypeStr.Needs ? CallType.Request : undefined),
    );

    this.isOffers$ = this.resourceType$.pipe(map(r => r === CallType.Offer));
    this.isNeeds$ = this.resourceType$.pipe(map(r => r === CallType.Request));
    this.accountId$ = this.route.parent.params.pipe(pluck('id'));

    this.resources$ = combineLatest([
      this.accountId$,
      this.resourceType$,
      this.route.data
    ]).pipe(
      switchMap(([accountId, resourceType, data]) => {
        if (data.accountType === AccountType.User) {
          return this.backend.searchResources(new SearchResourceRequest(undefined, resourceType, undefined, accountId, undefined, 10, 0));
        } else {
          return this.backend.searchResources(new SearchResourceRequest(undefined, resourceType, undefined, undefined, accountId, 10, 0));
        }
      }),
      pluck('resources')
    );
  }

}
