import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Observable, ReplaySubject, Subscription} from 'rxjs';
import {BackendService} from '../../api/backend.service';
import {distinctUntilChanged, filter, pluck, shareReplay, switchMap} from 'rxjs/operators';
import {GetGroupRequest, Group} from '../../api/models';

@Component({
  selector: 'app-group',
  template: `<ng-content></ng-content>`,
})
export class GroupComponent implements OnInit, OnDestroy {

  constructor(private backend: BackendService) {
  }

  idSubject = new ReplaySubject<string>();
  group$: Observable<Group> = this.idSubject.pipe(
    filter((id) => !!id),
    distinctUntilChanged(),
    switchMap(id => this.backend.getGroup(new GetGroupRequest(id))),
    pluck('group'),
    shareReplay()
  );

  private groupSub: Subscription;
  public group: Group = undefined;

  @Input()
  set id(value: string) {
    this.idSubject.next(value);
  }

  ngOnInit(): void {
    this.groupSub = this.group$.subscribe((g) => {
      this.group = g;
    });
  }

  ngOnDestroy(): void {
    this.groupSub.unsubscribe();
  }

}
