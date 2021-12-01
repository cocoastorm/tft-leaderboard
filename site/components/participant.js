import { indicator } from 'ordinal';
import styles from './participant.module.css';
import Tag from '../components/tag';

export default function ParticipantItem(props) {
  const profileIconName = `${props.summoner.profileIconId}.png`;

  return (
    <div className={styles.box}>
      <figure className="col-span-2 max-w-full rounded-xl sm:p-0 mx-auto">
        <img
          className="w-8 h-8 md:w-16 md:h-auto rounded-full mx-auto"
          src={'http://ddragon.leagueoflegends.com/cdn/11.23.1/img/profileicon/' + profileIconName}
        />
        <figcaption className="font-medium">
          <div className="text-center text-sm">{props.summoner.name}</div>
        </figcaption>
      </figure>

      { props.rank !== null &&
        <div>
          <span>{ props.rank.tier }</span>
          &nbsp;
          <span>{ props.rank.rank }</span>
        </div>
      }

      { props.rank !== null &&
        <div>
          <Tag value={ props.rank && props.rank.leaguePoints ? props.rank.leaguePoints : 0 } title="LP" />
        </div>
      }

      { props.rank !== null &&
        <div>
          <Tag value={ props.rank.wins } title="Wins" />
          <Tag value={ props.rank.losses } title="Losses" />
        </div>
      }

      { props.rank !== null && props.placement !== null &&
        <div class="ordinal ml-10">
          {props.placement}{indicator(props.placement)}
        </div>
      }
    </div>
  )
}
