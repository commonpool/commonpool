import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Action} from '../../api/models';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-offer-action',
  template: `
    <form class="d-inline-block" method="post" (ngSubmit)="submit()">
      <button
        type="submit"
        class="btn btn-sm"
        [disabled]="!action.enabled"
        [ngClass]="{
            'btn-danger': action.enabled && !action.completed && action.style === 'danger',
            'btn-outline-danger': action.completed && action.style === 'danger',
            'btn-success': action.enabled && !action.completed && action.style === 'success',
            'btn-outline-success': action.completed && action.style === 'success',
            'btn-secondary' : !action.enabled && !action.completed
          }">
        <app-check *ngIf="action.completed"></app-check>
        {{action.name}}
      </button>
    </form>`,
})
export class OfferActionComponent implements OnInit {

  constructor(private http: HttpClient) {
  }

  @Input()
  action: Action;

  @Output()
  submitted: EventEmitter<null> = new EventEmitter<null>();

  ngOnInit(): void {
  }

  submit() {
    this.http.post(this.action.actionUrl, null).subscribe(r => {
      this.submitted.next();
    });
  }
}
