import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';
import {BlockService} from '../block.service';

@Component({
  selector: 'app-block',
  styles: [`
    * {
      font-size: 95%;
    }
  `],
  template: `
    <ng-container [ngSwitch]="block.type">

      <div *ngSwitchCase="'section'" class="block block-section">
        <app-section-block [block]="block"></app-section-block>
      </div>

      <div *ngSwitchCase="'context'" class="block block-context">
        <app-context-block [block]="block"></app-context-block>
      </div>

      <div *ngSwitchCase="'actions'" class="block block-actions">
        <app-actions-block [block]="block"></app-actions-block>
      </div>

      <div *ngSwitchCase="'divider'" class="block block-divider">
        <app-divider-block></app-divider-block>
      </div>

      <div *ngSwitchCase="'image'" class="block block-images">
        <app-image-block [block]="block"></app-image-block>
      </div>

      <h5 *ngSwitchCase="'header'" class="block block-header">
        <app-header-block [block]="block"></app-header-block>
      </h5>

    </ng-container>
  `,
  providers: [BlockService]
})
export class BlockComponent {

  constructor(private blockSvc: BlockService) {
  }

  private _block: Block;

  @Input()
  set block(value: Block) {
    this._block = value;
    this.blockSvc.setBlock(value);
  }

  get block(): Block {
    return this._block;
  }

}
