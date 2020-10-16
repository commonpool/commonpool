import {Component} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {BehaviorSubject, combineLatest, of, ReplaySubject, Subject} from 'rxjs';
import {catchError, debounceTime, distinctUntilChanged, filter, map, shareReplay, startWith, switchMap, tap} from 'rxjs/operators';
import {Resource, ResourceType, SearchResourceRequest} from '../../api/models';

@Component({
  selector: 'app-resource-list-view',
  templateUrl: './resource-list-view.component.html',
  styleUrls: ['./resource-list-view.component.css']
})
export class ResourceListViewComponent {

  querySubject = new ReplaySubject<string>();
  query$ = this.querySubject.asObservable();

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
      filter(([q, t]) => !!q && q.length > 4),
      map(([q, t]) => new SearchResourceRequest(q, t, 10, 0)),
      tap(() => {
        this.pendingSubject.next(true);
      }),
      switchMap(req => {
        return this.backend.searchResources(req).pipe(
          map(r => r.resources),
          catchError(err => of([]))
        );
      }),
      tap(() => {
        this.pendingSubject.next(false);
      }),
      startWith([] as Resource[])
    );

  constructor(private backend: BackendService) {

  }


}
