import {Component, Input, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {Subject} from 'rxjs';
import {filter, pluck, switchMap, tap} from 'rxjs/operators';

@Component({
  selector: 'app-resource-name',
  templateUrl: './resource-name.component.html',
  styleUrls: ['./resource-name.component.css']
})
export class ResourceNameComponent implements OnInit {

  constructor(private backend: BackendService) {
  }

  resourceName: string;
  resourceIdSubject = new Subject<string>();
  resourceName$ = this.resourceIdSubject.asObservable().pipe(
    filter(id => !!id),
    switchMap(id => this.backend.getResource(id)),
    pluck('resource', 'info', 'name')
  ).subscribe((a) => {
    console.log(a);
    this.resourceName = a;
  });

  ngOnInit(): void {
  }

  @Input()
  set resourceId(value: string) {
    this.resourceIdSubject.next(value);
  }

}
