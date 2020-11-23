import {Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
import {
  ActionsBlock,
  Attachment,
  Block,
  ButtonElement,
  ButtonStyle,
  ContextBlock,
  HeaderBlock,
  ImageElement, Message,
  SectionBlock,
  TextObject,
  TextType
} from '../../api/models';

interface Payload {
  blocks: Block[];
  attachments: Attachment[];
}

@Component({
  selector: 'app-blocks',
  templateUrl: './blocks.component.html',
  styleUrls: ['./blocks.component.css'],
  encapsulation: ViewEncapsulation.None
})
export class BlocksComponent implements OnInit {

  @Input()
  public message: Message;

  constructor() {
  }

  ngOnInit(): void {
  }

}
