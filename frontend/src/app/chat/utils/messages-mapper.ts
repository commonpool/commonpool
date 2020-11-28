import {Message} from '../../api/models';

export class MessageGroup {
  public constructor(public sentByUsername: string, public sentById: string, public formattedDate: string, public messages: Message[]) {
  }
}

export function MapMessages(messages: Message[]): MessageGroup[] {

  if (!messages || messages.length === 0) {
    return [];
  }

  const result: MessageGroup[] = [];
  const sortedMessages = messages.sort((a, b) => a.sentAtDate.valueOf() - b.sentAtDate.valueOf());

  let previousFormattedDate: string;
  let previousMessage: Message;
  let previousSentBy: string;
  for (const message of sortedMessages) {

    const messageDateStr = getMessageGroupFormattedDate(message);
    let isNewGroup = false;

    if (previousFormattedDate !== messageDateStr) {
      result.push(new MessageGroup(message.sentByUsername, message.sentById, messageDateStr, []));
      isNewGroup = true;
    }

    if (!isNewGroup && previousSentBy !== message.sentById) {
      result.push(new MessageGroup(message.sentByUsername, message.sentById, messageDateStr, []));
      isNewGroup = true;
    }

    result[result.length - 1].messages.push(message);

    previousMessage = message;
    previousSentBy = message.sentById;
    previousFormattedDate = messageDateStr;

  }

  return result;

}

export function getMessageGroupFormattedDate(message: Message): string {

  const messageDate = message.sentAtDate;
  const now = new Date();
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);

  if (isSameDayAs(now, messageDate)) {
    return 'today';
  } else if (isSameDayAs(yesterday, messageDate)) {
    return 'yesterday';
  } else {
    return messageDate.toDateString();
  }

}

export function isSameDayAs(date1: Date, date2: Date) {
  return date1.getFullYear() === date2.getFullYear() && date1.getMonth() === date2.getMonth() && date1.getDate() === date2.getDate();
}

