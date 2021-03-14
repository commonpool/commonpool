import {Component, OnDestroy, OnInit} from '@angular/core';
import {ResourceEvaluationService} from './resource-evaluation.service';

@Component({
  template: `
    <div *ngIf="show" class="modal-backdrop fade show" (click)="close()"></div>
    <div *ngIf="show" (click)="doNothing($event)" class="modal modal-open fade show modal-close" tabindex="-1"
         style="display:block">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Modal title</h5>
            <button type="button" class="btn-close" (click)="close()"></button>
          </div>
          <div class="modal-body">
            ...
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary modal-close">Close</button>
            <button type="button" class="btn btn-primary modal-close">Save changes</button>
          </div>
        </div>
      </div>
    </div>
    <button (click)="triggerEvaluation()">Show modal</button>
  `,
})
export class ResourceEvaluationComponent implements OnInit, OnDestroy {

  show = false;

  public constructor(private svc: ResourceEvaluationService) {
  }

  sub = this.svc.resourceEvaluation$.subscribe((id) => {
    this.show = true;
  });

  close() {
    this.show = false;
  }

  open() {
    this.show = true;
  }

  triggerEvaluation() {
    this.svc.evaluateResource('');
  }

  doNothing($event: MouseEvent) {
    if (($event.target as HTMLElement).classList.contains('modal-close')) {
      this.close();
    }
    $event.preventDefault();
    $event.stopPropagation();
  }

  ngOnInit(): void {

  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

}
