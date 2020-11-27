import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-actions-block',
  template: `
    <div class="d-flex">
      <ng-container *ngFor="let blockElement of block.elements; let i = index">
        <ng-container *ngIf="blockElement.type === 'button'">
          <div class="actions-container">
            <app-button [buttonElement]="blockElement"></app-button>
          </div>
        </ng-container>
      </ng-container>
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
