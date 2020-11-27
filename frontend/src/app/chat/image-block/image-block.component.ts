import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-image-block',
  template: `
    <ng-container *ngIf="block.title">
      <div class="my-1">
        <app-text-object [textObject]="block.title" [subtle]="true" [small]="true"></app-text-object>
      </div>
    </ng-container>
    <img [src]="block.imageUrl" style="max-height: 16rem" class="rounded border shadow-sm">`
})
export class ImageBlockComponent {

  constructor() {
  }

  @Input()
  block: Block;

}
