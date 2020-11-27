import {Component, Input, OnInit} from '@angular/core';
import {Block} from '../../api/models';
import {BlockService} from '../block.service';

@Component({
  selector: 'app-block',
  templateUrl: './block.component.html',
  styleUrls: ['./block.component.css'],
  providers: [BlockService]
})
export class BlockComponent implements OnInit {

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

  ngOnInit(): void {
  }

}
