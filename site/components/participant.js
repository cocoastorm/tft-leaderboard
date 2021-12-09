import { indicator } from 'ordinal';
import Tag from './tag';

export default function ParticipantItem(props) {
  const profileIconName = `${props.summoner.profileIconId}.png`;
  
  let lastItemClassName = 'mx-5';
  if (props.rank === null) {
    lastItemClassName = `${lastItemClassName} flex-1`;
  }

  return (
    <div className="bg-gradient-to-r from-sky-700 hover:from-sky-600 to-sky-800 hover:to-sky-700 rounded-md max-w-2xl py-4 mb-5 md:w-1/2">
      <div className="flex flex-1 items-center justify-around w-2/3 md:w-full">
        { props.placement !== null &&
          <div className="ordinal mx-5">
            {props.placement}{indicator(props.placement)}
          </div>
        }

        <div className="max-w-full rounded-xl mx-2 p-8 sm:p-0">
          <img
            className="w-8 h-8 md:w-16 md:h-auto rounded-full mx-auto"
            src={'http://ddragon.leagueoflegends.com/cdn/11.23.1/img/profileicon/' + profileIconName}
          />
        </div>

        <div className="flex flex-col content-center w-3/5 md:w-1/5 mx-3">
          <div className="font-medium text-center text-sm">{props.summoner.name}</div>

          { props.rank !== null &&
            <div className="text-center font-semibold text-sm text-teal-400">
              <span className="pr-1">{ props.rank.tier }</span>
              <span>{ props.rank.rank }</span>
            </div>
          }
        </div>

        <div className="flex flex-col items-center mx-3">
          <Tag className="text-emerald-400" value={ (props.rank && props.rank.wins) ? props.rank.wins : 0 } title="Wins" />
          <Tag className="text-fuchsia-400" value={ (props.rank && props.rank.losses) ? props.rank.losses : 0 } title="Losses" />
        </div>

        <div className="text-teal-400">
          <Tag value={ props.rank && props.rank.leaguePoints ? props.rank.leaguePoints : 0 } title="LP" />
        </div>
      </div>
    </div>
  )
}
