import {Injectable} from '@angular/core';
import {Message, SubmitAction, SubmitInteractionPayload, SubmitInteractionRequest} from '../api/models';
import {BackendService} from '../api/backend.service';

@Injectable()
export class BlocksService {
  constructor(private backend: BackendService) {
  }

  private _message: Message;

  public getMessage(): Message {
    return this._message;
  }

  public setMessage(message: Message) {
    this._message = message;
  }

  public submitInteraction(action: SubmitAction) {
    this.backend.submitMessageInteraction(new SubmitInteractionRequest(new SubmitInteractionPayload(
      this._message.id,
      [action],
      {}
    ))).subscribe((res) => {
      console.log(res);
    });
  }

}
