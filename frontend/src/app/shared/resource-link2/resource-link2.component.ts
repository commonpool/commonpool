import {Component, Input, OnInit} from '@angular/core';
import {ReplaySubject} from 'rxjs';
import {BackendService} from '../../api/backend.service';
import {distinctUntilChanged, filter, shareReplay, switchMap} from 'rxjs/operators';

@Component({
  selector: 'app-resource-link2',
  templateUrl: './resource-link2.component.html',
  styleUrls: ['./resource-link2.component.css']
})
export class ResourceLink2Component implements OnInit {

  constructor(private backend: BackendService) {
  }

  idSubject = new ReplaySubject<string>();
  resource$ = this.idSubject.pipe(
    filter((id) => !!id),
    distinctUntilChanged(),
    switchMap(id => this.backend.getResource(id)),
    shareReplay()
  );

  @Input()
  set id(value: string) {
    this.idSubject.next(value);
  }

  ngOnInit(): void {
  }

}
