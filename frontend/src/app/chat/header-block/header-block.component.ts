import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-header-block',
  template: `
    <div class="font-weight-bold block-header mt-2">
      <app-text-object [textObject]="block" [small]="true"></app-text-object>
    </div>
  `
})
export class HeaderBlockComponent {

  constructor() {
  }

  @Input()
  block: Block;

}
