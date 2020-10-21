import {Component} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {combineLatest, of, ReplaySubject} from 'rxjs';
import {catchError, debounceTime, distinctUntilChanged, map, startWith, switchMap, tap} from 'rxjs/operators';
import {ExtendedResource, ResourceType, SearchResourceRequest} from '../../api/models';
import {AuthService} from '../../auth.service';

@Component({
  selector: 'app-resource-list-view',
  templateUrl: './resource-list-view.component.html',
  styleUrls: ['./resource-list-view.component.css']
})
export class ResourceListViewComponent {

  querySubject = new ReplaySubject<string>();
  query$ = this.querySubject.asObservable().pipe(startWith(''));

  pendingSubject = new ReplaySubject<boolean>();
  pending$ = this.pendingSubject.asObservable().pipe(startWith(false));

  resourceTypeSubject = new ReplaySubject<ResourceType>();
  resourceType$ = this.resourceTypeSubject.asObservable().pipe(startWith(ResourceType.Offer));

  resources$ =
    combineLatest(
      [
        this.query$.pipe(
          debounceTime(100),
          distinctUntilChanged()
        ),
        this.resourceType$
      ]
    ).pipe(
      map(([q, t]) => new SearchResourceRequest(q, t, undefined, 10, 0)),
      tap(() => {
        this.pendingSubject.next(true);
      }),
      switchMap(req => {
        return this.backend.searchResources(req).pipe(
          map(r => r.resources),
          catchError(err => of([] as ExtendedResource[]))
        );
      }),
      tap(() => {
        this.pendingSubject.next(false);
      }),
      startWith([] as ExtendedResource[])
    );

  constructor(private backend: BackendService) {

  }


}
