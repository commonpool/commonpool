import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-section-block',
  template: `
    <div class="d-flex flex-row">

      <ng-container *ngIf="block.text">
        <app-text-object [textObject]="block.text"></app-text-object>
      </ng-container>

      <div class="flex-grow-1" style="max-width: 10rem"></div>

      <ng-container *ngIf="block.accessory; let accessory">
        <ng-container *ngIf="accessory.type === 'button'">
          <div style="position:relative; top:-0.25rem">
            <app-button [buttonElement]="accessory"></app-button>
          </div>
        </ng-container>
      </ng-container>

    </div>
  `,
})
export class SectionBlockComponent {

  constructor() {
  }

  @Input()
  block: Block;

}
