import {Component} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {ActivatedRoute} from '@angular/router';
import {pluck, switchMap} from 'rxjs/operators';
import {Observable, of} from 'rxjs';
import {GroupReportResponse} from '../../api/models';

@Component({
  template: `
    <ng-container *ngIf="history$ | async; let history">
      <table class="table">
        <thead>
        <tr>
          <th>Activity</th>
          <th>Items Received</th>
          <th>Items Given</th>
          <th>Items Owned</th>
          <th>Items Lent</th>
          <th>Items Borrowed</th>
          <th>Services Given</th>
          <th>Services Received</th>
          <th>Offers</th>
          <th>Requests</th>
          <th>Hours in Bank</th>
        </tr>
        </thead>
        <tbody>
        <tr *ngFor="let item of history.entries">
          <td>{{item.activity}}</td>
          <td>{{item.itemsReceived}}</td>
          <td>{{item.itemsGiven}}</td>
          <td>{{item.itemsOwned}}</td>
          <td>{{item.itemsLent}}</td>
          <td>{{item.itemsBorrowed}}</td>
          <td>{{item.servicesGiven}}</td>
          <td>{{item.servicesReceived}}</td>
          <td>{{item.offerCount}}</td>
          <td>{{item.requestCount}}</td>
          <td>
            <app-duration [duration]="item.hoursInBank"></app-duration>
          </td>
          <td></td>
        </tr>
        </tbody>
      </table>
    </ng-container>
  `
})
export class GroupHistoryComponent {

  public constructor(private backend: BackendService, private route: ActivatedRoute) {
  }

  groupId$ = this.route.parent.params.pipe(pluck('id'));

  history$: Observable<GroupReportResponse> = this.groupId$.pipe(
    switchMap(value => {
      return this.backend.groupHistoryReport$(of(value));
    })
  );

}
