import {Injectable} from '@angular/core';
import {Block} from '../api/models';

@Injectable()
export class BlockService {

  constructor() {
  }

  private _block: Block;

  public getBlock(): Block {
    return this._block;
  }

  public setBlock(block: Block) {
    this._block = block;
  }

}
