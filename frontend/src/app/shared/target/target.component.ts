import {Component, Input, OnInit} from '@angular/core';
import {ReplaySubject, Subject} from 'rxjs';
import {Target} from '../../api/models';

@Component({
  selector: 'app-target',
  templateUrl: './target.component.html',
  styleUrls: ['./target.component.css']
})
export class TargetComponent implements OnInit {

  constructor() {
  }

  private targetSubject = new ReplaySubject<Target>(1);
  public target$ = this.targetSubject.asObservable();

  @Input()
  set target(value: Target) {
    this.targetSubject.next(value);
  }



  ngOnInit(): void {
  }

}
