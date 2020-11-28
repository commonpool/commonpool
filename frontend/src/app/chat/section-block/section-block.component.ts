import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-section-block',
  styles: [`
    .section-row {
      align-items: center;
    }
  `],
  template: `
    <div class="section-row d-flex flex-row">

      <ng-container *ngIf="block.text">
        <app-text-object [textObject]="block.text"></app-text-object>
      </ng-container>

      <div class="flex-grow-1" style="max-width: 10rem"></div>

      <ng-container *ngIf="block.accessory; let accessory">
        <ng-container *ngIf="accessory.type === 'button'">
          <app-button [buttonElement]="accessory"></app-button>
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
