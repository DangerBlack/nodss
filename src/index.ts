import * as TelegramBot from 'node-telegram-bot-api';
import {MessageEntity, SendMessageOptions} from 'node-telegram-bot-api';
import {ReplyData, ReplyManager} from 'node-telegram-operation-manager';
import { create_pie } from './chart_adapter';
import { Dss } from './dss';
const Joi = require('joi');
require('dotenv').config();

console.log('Configuring the token');
const token = process.env.TOKEN;

console.log(`Token ${token.substr(0,3)}...${token.substr(token.length-3)}`);

const bot = new TelegramBot(token, {polling: true});
const reply = new ReplyManager();
const parse_mode: SendMessageOptions = {parse_mode: 'HTML'};

bot.onText(/\/pie/, async (msg) => 
{
  const chat_id = msg.from.id;
  bot.sendMessage(chat_id, 'Please send the statement');

  reply.register(chat_id, (data?: ReplyData) => 
  {
    const head = data.text;

    bot.sendMessage(chat_id, `Send me the first option for (<i>"${head}"</i>)`, parse_mode);

    return {
      repeat: false,
      head,
      values: []
    };
  }).register(chat_id, (data?: ReplyData) =>
  {
    const head = data.previousData.head;
    const repeat = data.text !== '/done';
    const values = data.previousData.values;

    console.log(repeat);
    if(repeat)
    {
      values.push(data.text);
      bot.sendMessage(chat_id, `Send me the next option for (<i>"${head}"</i>)\nSend /done when you have done.`, parse_mode);
    }
    else
    {
      const samples = values.map(value => `- ${head} ${value}`).join('\n');

      bot.sendMessage(chat_id, `Elaborating plot for <i>"${head}"</i>\n${samples}`, parse_mode);
      bot.sendChatAction(chat_id, 'upload_photo');
      Dss.count_events(head, values).then(data => 
      {
        create_pie(head, values, data.values).then(img => bot.sendPhoto(chat_id, img, {}));
      }).catch((_error)=>
      {
        bot.sendMessage(chat_id, `Unable to complete the operation for <i>"${head}"</i>.`, parse_mode);
      });
    }

    return {
      repeat,
      head,
      values
    };
  });
});

bot.onText(/\/start/, (msg) => 
{
  const chatId = msg.from.id;
  bot.sendMessage(chatId, `Hello, I\'m a @dsspiebot!
Send me a phrase and some different conclusion of that phrase and I tell you which is the online popularity.
Example:
/pie
   I like
   foods
   cats`);
});

bot.on('message', (msg) => 
{
  if((command_exclusion(msg.text) || !has_entity('bot_command', msg.entities)) && reply.expects(msg.from.id)) 
  {
		const { text, entities } = msg;
		reply.execute(msg.from.id, { text, entities });
	}
});

function command_exclusion(command: string)
{
  switch(command)
  {
    case '/done':
      return true;
    default:
      return false;
  }
}
function has_entity(entity: string, entities?: MessageEntity[]) 
{
	if(!entities || !entities.length) 
	  return false;

	return entities.some(e => e.type === entity);
}