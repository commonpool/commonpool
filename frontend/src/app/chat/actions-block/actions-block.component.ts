import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-actions-block',
  template: `
    <div class="d-flex">
      <div class="actions-container">
        <ng-container *ngFor="let blockElement of block.elements; let i = index">
          <ng-container *ngIf="blockElement.type === 'button'">
            <app-button [buttonElement]="blockElement"></app-button>
          </ng-container>
        </ng-container>
      </div>
    </div>
  `,
  styles: [`
    .actions-container > :not(:first-child) {
      margin-left: 0.125rem;
    }`
  ],
})
export class ActionsBlockComponent {

  constructor() {
  }

  @Input()
  block: Block;

}
