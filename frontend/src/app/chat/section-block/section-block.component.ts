import {ChangeDetectionStrategy, Component, Input} from '@angular/core';
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
        <div class="flex-grow-1">
            <app-text-object class="w-100" [textObject]="block.text"></app-text-object>
        </div>
      </ng-container>

      <div class="flex-grow-1" style="max-width: 10rem"></div>

      <ng-container *ngIf="block.accessory; let accessory">
        <ng-container *ngIf="accessory.type === 'button'">
          <div class="flex-grow-1">
            <app-button [buttonElement]="accessory"></app-button>
          </div>
        </ng-container>
      </ng-container>

    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SectionBlockComponent {

  constructor() {
    console.log("new section block")
  }

  @Input()
  block: Block;

}
